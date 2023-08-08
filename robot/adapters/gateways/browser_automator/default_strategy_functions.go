package browser_automator

import (
	"automator-go/robot/entities/models"
	"automator-go/robot/entities/validation"
	"automator-go/robot/usecases/task"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05Z04:00"
)

func findElement(page *rod.Page, action models.TaskAction) (*rod.Element, error) {
	var element *rod.Element
	var err error
	if validation.IsXpath(action.Value) {
		element, err = page.ElementX(action.Value)
	} else {
		element, err = page.Element(action.Value)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting element by selector: %w", err)
	}

	return element, nil
}

func click(page *rod.Page, action models.TaskAction) error {
	element, err := findElement(page, action)
	if err != nil {
		return err
	}

	err = element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("error clicking element: %w", err)
	}

	return nil
}

func navigate(page *rod.Page, action models.TaskAction) error {
	err := click(page, action)
	if err != nil {
		return err
	}

	page.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)

	return nil
}

func scrollDown(page *rod.Page, action models.TaskAction) error {
	parsedSteps, err := strconv.ParseInt(action.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing scroll steps: %w", err)
	}

	window, err := page.GetWindow()
	if err != nil {
		return fmt.Errorf("error getting page window: %w", err)
	}

	currentPos := float64(*window.Top)
	for i := 0; i < int(parsedSteps); i++ {
		currentPos += 100
		err = page.Mouse.Scroll(0, currentPos, 0)
		if err != nil {
			return fmt.Errorf("error scrolling down: %w", err)
		}
	}

	return nil
}

func captureElement(element *rod.Element) ([]byte, error) {
	bin, err := element.Screenshot(proto.PageCaptureScreenshotFormatPng, 0)

	if err != nil {
		return nil, fmt.Errorf("error capturing element: %w", err)
	}

	return bin, nil
}

func _captureAction(page *rod.Page, element *rod.Element, resource bool) (*task.RawMedia, error) {
	mediaScreenshot, err := captureElement(element)
	if err != nil {
		return nil, err
	}

	pageScreenshot, err := page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
	})
	if err != nil {
		return nil, fmt.Errorf("error capturing page: %w", err)
	}

	shape, err := element.Shape()
	if err != nil {
		return nil, fmt.Errorf("error getting element shape: %w", err)
	}

	box := shape.Box()

	info, err := page.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting page info: %w", err)
	}

	extension := "png"

	var resourceCapture []byte
	if resource {
		err = element.WaitLoad()
		if err != nil {
			return nil, fmt.Errorf("error waiting resource to load: %w", err)
		}

		src, err := element.Property("currentSrc")
		if err != nil {
			return nil, fmt.Errorf("error getting element currentSrc: %w", err)
		}
		srcSplit := strings.Split(src.String(), ".")
		extension = srcSplit[len(srcSplit)-1]

		video := false
		for _, videoExtension := range VideoExtensions {
			if extension == videoExtension {
				video = true
				break
			}
		}

		if video {
			httpClient := http.DefaultClient
			getVideo, err := httpClient.Get(src.String())
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					println("error closing video body")
				}
			}(getVideo.Body)
			if err != nil {
				return nil, fmt.Errorf("error getting video: %w", err)
			}

			resourceCapture, err = io.ReadAll(getVideo.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading video: %w", err)
			}
		} else {
			resourceCapture, err = element.Resource()
			if err != nil {
				return nil, fmt.Errorf("error getting element resource: %w", err)
			}
		}
	}

	return &task.RawMedia{
		Ext:        extension,
		Media:      mediaScreenshot,
		Screenshot: pageScreenshot,
		Resource:   resourceCapture,
		Height:     box.Height,
		Width:      box.Width,
		X:          box.X,
		Y:          box.Y,
		Url:        info.URL,
	}, nil
}

func capture(page *rod.Page, action models.TaskAction) (*task.RawMedia, error) {
	element, err := findElement(page, action)
	if err != nil {
		return nil, err
	}

	timeOutStableEnv := os.Getenv("BROWSER_WAIT_STABLE_TIMEOUT")
	if strings.TrimSpace(timeOutStableEnv) == "" {
		timeOutStableEnv = "5s"
	}

	timeOutStable, err := time.ParseDuration(timeOutStableEnv)
	if err != nil {
		return nil, fmt.Errorf("error parsing timeout stable env: %w", err)
	}

	if err = page.WaitStable(timeOutStable); err != nil {
		return nil, fmt.Errorf("error waiting element to load: %w", err)
	}

	return _captureAction(page, element, false)
}

func waitSeconds(action models.TaskAction) error {
	parsedSeconds, err := strconv.ParseInt(action.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing seconds: %w", err)
	}

	seconds := int(parsedSeconds)
	if seconds < 1 {
		return fmt.Errorf("seconds must be greater than 0")
	}

	time.Sleep(time.Duration(seconds) * time.Second)

	return nil
}

func writeInput(page *rod.Page, action models.TaskAction) error {
	element, err := findElement(page, action)
	if err != nil {
		return err
	}

	err = element.Input(action.Value)
	if err != nil {
		return fmt.Errorf("error writing input: %w", err)
	}

	return nil
}

func clearInput(page *rod.Page, action models.TaskAction) error {
	element, err := findElement(page, action)
	if err != nil {
		return err
	}

	err = element.SelectAllText()
	if err != nil {
		return fmt.Errorf("error selecting text on input: %w", err)
	}

	err = element.Input("")
	if err != nil {
		return fmt.Errorf("error clearing input: %w", err)
	}

	return nil
}

func selectOptions(page *rod.Page, action models.TaskAction) error {
	element, err := findElement(page, action)
	if err != nil {
		return err
	}

	parsedOptions := strings.Split(action.Value, ",")

	err = element.Select(parsedOptions, true, rod.SelectorTypeText)
	if err != nil {
		return fmt.Errorf("error selecting option: %w", err)
	}

	return nil
}

func writeTime(page *rod.Page, action models.TaskAction) error {
	element, err := findElement(page, action)
	if err != nil {
		return err
	}

	timeToWrite, err := time.Parse(timeLayout, action.Value)

	err = element.InputTime(timeToWrite)
	if err != nil {
		return fmt.Errorf("error writing time: %w", err)
	}

	return nil
}

func downloadResource(page *rod.Page, action models.TaskAction) (*task.RawMedia, error) {
	element, err := findElement(page, action)
	if err != nil {
		return nil, err
	}

	return _captureAction(page, element, true)
}
