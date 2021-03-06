package handlers

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

// CreateEditTag creates a new tag on current tenant
func CreateEditTag() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.CreateEditTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		var (
			tag *models.Tag
			err error
		)

		if input.Model.Slug != "" {
			tag, err = c.Services().Tags.Update(input.Tag.ID, input.Model.Name, input.Model.Color, input.Model.IsPublic)
		} else {
			tag, err = c.Services().Tags.Add(input.Model.Name, input.Model.Color, input.Model.IsPublic)
		}

		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(tag)
	}
}

// RemoveTag deletes anexisting tag
func RemoveTag() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.RemoveTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tags.Remove(input.Tag.ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// AssignTag to existing dea
func AssignTag() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.AssignUnassignTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tags.AssignTag(input.Tag.ID, input.Idea.ID, c.User().ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// UnassignTag from existing dea
func UnassignTag() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.AssignUnassignTag)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := c.Services().Tags.UnassignTag(input.Tag.ID, input.Idea.ID)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}
