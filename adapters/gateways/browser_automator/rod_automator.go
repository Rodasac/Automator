package browser_automator

import (
	"automator-go/entities/models"
	"automator-go/usecases/task"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
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
			err = navigateXpath(page, action)
			if err != nil {
				return nil, fmt.Errorf("error navigating xpath: %w", err)
			}
		case models.NavigateSelector:
			err = navigateSelector(page, action)
			if err != nil {
				return nil, fmt.Errorf("error navigating css: %w", err)
			}
		case models.ClickXpath:
			err = clickXpath(page, action)
			if err != nil {
				return nil, fmt.Errorf("error clicking xpath: %w", err)
			}
		case models.ClickSelector:
			err = clickSelector(page, action)
			if err != nil {
				return nil, fmt.Errorf("error clicking css: %w", err)
			}
		case models.ScrollDown:
			err = scrollDown(page, action)
			if err != nil {
				return nil, fmt.Errorf("error scrolling down: %w", err)
			}
		case models.CaptureXpath:
			rawMedia, err := captureXpath(page, action)
			if err != nil {
				return nil, fmt.Errorf("error capturing xpath: %w", err)
			}
			rawMedias = append(rawMedias, *rawMedia)
		case models.CaptureSelector:
			rawMedia, err := captureSelector(page, action)
			if err != nil {
				return nil, fmt.Errorf("error capturing css: %w", err)
			}
			rawMedias = append(rawMedias, *rawMedia)
		case models.WaitSeconds:
			err = waitSeconds(action)
			if err != nil {
				return nil, fmt.Errorf("error waiting seconds: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown action type: %s", action.Type.String())
		}
	}

	return &rawMedias, nil
}
