package browser_automator

import (
	"automator-go/entities/models"
	"automator-go/usecases/task"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"strconv"
	"time"
)

type RodAutomator struct {
	browser *rod.Browser
}

func NewRodAutomator(browser *rod.Browser) *RodAutomator {
	return &RodAutomator{browser: browser}
}

func (at *RodAutomator) Run(taskToRun *models.Task) (*[]task.RawMedia, error) {
	err := at.browser.Connect()
	if err != nil {
		return nil, fmt.Errorf("error connecting to browser: %w", err)
	}

	page, err := at.browser.Page(proto.TargetCreateTarget{URL: taskToRun.Url})
	if err != nil {
		return nil, fmt.Errorf("error creating page: %w", err)
	}
	defer func(page *rod.Page) {
		err := page.Close()
		if err != nil {
			println("error closing page: " + err.Error())
		}
	}(page)

	err = page.WaitStable(800*time.Millisecond, 1)
	if err != nil {
		return nil, fmt.Errorf("error waiting for page to be stable: %w", err)
	}

	rawMedias := make([]task.RawMedia, 0)

	for _, action := range taskToRun.Actions {
		switch action.Type {
		case models.NavigateXpath:
			err = at.navigateXpath(page, action)
			if err != nil {
				return nil, fmt.Errorf("error navigating xpath: %w", err)
			}
		case models.NavigateSelector:
			err = at.navigateSelector(page, action)
			if err != nil {
				return nil, fmt.Errorf("error navigating css: %w", err)
			}
		case models.ClickXpath:
			err = at.clickXpath(page, action)
			if err != nil {
				return nil, fmt.Errorf("error clicking xpath: %w", err)
			}
		case models.ClickSelector:
			err = at.clickSelector(page, action)
			if err != nil {
				return nil, fmt.Errorf("error clicking css: %w", err)
			}
		case models.ScrollDown:
			err = at.scrollDown(page, action)
			if err != nil {
				return nil, fmt.Errorf("error scrolling down: %w", err)
			}
		case models.CaptureXpath:
			rawMedia, err := at.captureXpath(page, action)
			if err != nil {
				return nil, fmt.Errorf("error capturing xpath: %w", err)
			}
			rawMedias = append(rawMedias, *rawMedia)
		case models.CaptureSelector:
			rawMedia, err := at.captureSelector(page, action)
			if err != nil {
				return nil, fmt.Errorf("error capturing css: %w", err)
			}
			rawMedias = append(rawMedias, *rawMedia)
		default:
			return nil, fmt.Errorf("unknown action type: %s", action.Type)
		}
	}

	return &rawMedias, nil
}

func (at *RodAutomator) findElementX(page *rod.Page, action models.TaskAction) (*rod.Element, error) {
	element, err := page.ElementX(action.Value)
	if err != nil {
		return nil, fmt.Errorf("error getting element by xpath: %w", err)
	}

	return element, nil
}

func (at *RodAutomator) findElementSelector(page *rod.Page, action models.TaskAction) (*rod.Element, error) {
	element, err := page.Element(action.Value)
	if err != nil {
		return nil, fmt.Errorf("error getting element by css selector: %w", err)
	}

	return element, nil
}

func (at *RodAutomator) click(element *rod.Element) error {
	err := element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("error clicking element: %w", err)
	}

	return nil
}

func (at *RodAutomator) navigateXpath(page *rod.Page, action models.TaskAction) error {
	element, err := at.findElementX(page, action)
	if err != nil {
		return err
	}

	err = at.click(element)
	if err != nil {
		return err
	}

	page.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)

	return nil
}

func (at *RodAutomator) navigateSelector(page *rod.Page, action models.TaskAction) error {
	element, err := at.findElementSelector(page, action)
	if err != nil {
		return err
	}

	err = at.click(element)
	if err != nil {
		return err
	}

	page.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)

	return nil
}

func (at *RodAutomator) clickXpath(page *rod.Page, action models.TaskAction) error {
	element, err := at.findElementX(page, action)
	if err != nil {
		return err
	}

	err = at.click(element)
	if err != nil {
		return err
	}

	return nil
}

func (at *RodAutomator) clickSelector(page *rod.Page, action models.TaskAction) error {
	element, err := at.findElementSelector(page, action)
	if err != nil {
		return err
	}

	err = at.click(element)
	if err != nil {
		return err
	}

	return nil
}

func (at *RodAutomator) scrollDown(page *rod.Page, action models.TaskAction) error {
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

func (at *RodAutomator) captureElement(element *rod.Element) ([]byte, error) {
	bin, err := element.Screenshot(proto.PageCaptureScreenshotFormatPng, 0)
	if err != nil {
		return nil, fmt.Errorf("error capturing element: %w", err)
	}

	return bin, nil
}

func (at *RodAutomator) capture(page *rod.Page, element *rod.Element) (*task.RawMedia, error) {
	mediaScreenshot, err := at.captureElement(element)
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

	return &task.RawMedia{
		Media:      mediaScreenshot,
		Screenshot: pageScreenshot,
		Height:     box.Height,
		Width:      box.Width,
		X:          box.X,
		Y:          box.Y,
		Url:        info.URL,
	}, nil
}

func (at *RodAutomator) captureXpath(page *rod.Page, action models.TaskAction) (*task.RawMedia, error) {
	element, err := at.findElementX(page, action)
	if err != nil {
		return nil, err
	}

	return at.capture(page, element)
}

func (at *RodAutomator) captureSelector(page *rod.Page, action models.TaskAction) (*task.RawMedia, error) {
	element, err := at.findElementSelector(page, action)
	if err != nil {
		return nil, err
	}

	return at.capture(page, element)
}
