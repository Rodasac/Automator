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

func findElementX(page *rod.Page, action models.TaskAction) (*rod.Element, error) {
	element, err := page.ElementX(action.Value)
	if err != nil {
		return nil, fmt.Errorf("error getting element by xpath: %w", err)
	}

	return element, nil
}

func findElementSelector(page *rod.Page, action models.TaskAction) (*rod.Element, error) {
	element, err := page.Element(action.Value)
	if err != nil {
		return nil, fmt.Errorf("error getting element by css selector: %w", err)
	}

	return element, nil
}

func click(element *rod.Element) error {
	err := element.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("error clicking element: %w", err)
	}

	return nil
}

func navigateXpath(page *rod.Page, action models.TaskAction) error {
	element, err := findElementX(page, action)
	if err != nil {
		return err
	}

	err = click(element)
	if err != nil {
		return err
	}

	page.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)

	return nil
}

func navigateSelector(page *rod.Page, action models.TaskAction) error {
	element, err := findElementSelector(page, action)
	if err != nil {
		return err
	}

	err = click(element)
	if err != nil {
		return err
	}

	page.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)

	return nil
}

func clickXpath(page *rod.Page, action models.TaskAction) error {
	element, err := findElementX(page, action)
	if err != nil {
		return err
	}

	err = click(element)
	if err != nil {
		return err
	}

	return nil
}

func clickSelector(page *rod.Page, action models.TaskAction) error {
	element, err := findElementSelector(page, action)
	if err != nil {
		return err
	}

	err = click(element)
	if err != nil {
		return err
	}

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

func capture(page *rod.Page, element *rod.Element) (*task.RawMedia, error) {
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

func captureXpath(page *rod.Page, action models.TaskAction) (*task.RawMedia, error) {
	element, err := findElementX(page, action)
	if err != nil {
		return nil, err
	}

	return capture(page, element)
}

func captureSelector(page *rod.Page, action models.TaskAction) (*task.RawMedia, error) {
	element, err := findElementSelector(page, action)
	if err != nil {
		return nil, err
	}

	return capture(page, element)
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
