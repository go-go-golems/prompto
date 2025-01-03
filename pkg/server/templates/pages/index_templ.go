// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/go-go-golems/prompto/pkg"
	"github.com/go-go-golems/prompto/pkg/server/templates"
	"github.com/go-go-golems/prompto/pkg/server/templates/components"
)

func copyToClipboard(text string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_copyToClipboard_c982`,
		Function: `function __templ_copyToClipboard_c982(text){copyToClipboard(text)
}`,
		Call:       templ.SafeScript(`__templ_copyToClipboard_c982`, text),
		CallInline: templ.SafeScriptInline(`__templ_copyToClipboard_c982`, text),
	}
}

func addToFavorites(name string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_addToFavorites_543a`,
		Function: `function __templ_addToFavorites_543a(name){addToFavorites(name)
}`,
		Call:       templ.SafeScript(`__templ_addToFavorites_543a`, name),
		CallInline: templ.SafeScriptInline(`__templ_addToFavorites_543a`, name),
	}
}

func Index(repositories []string, repos map[string]*pkg.Repository) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script src=\"/static/js/favorites.js\"></script> <div class=\"toast-container position-fixed bottom-0 end-0 p-3\"><div id=\"copyToast\" class=\"toast align-items-center text-bg-success\" role=\"alert\" aria-live=\"assertive\" aria-atomic=\"true\"><div class=\"d-flex\"><div class=\"toast-body\"><i class=\"bi bi-clipboard-check me-2\"></i>Prompt copied to clipboard!</div><button type=\"button\" class=\"btn-close btn-close-white me-2 m-auto\" data-bs-dismiss=\"toast\" aria-label=\"Close\"></button></div></div><div id=\"favToast\" class=\"toast align-items-center text-bg-primary\" role=\"alert\" aria-live=\"assertive\" aria-atomic=\"true\"><div class=\"d-flex\"><div class=\"toast-body\"><i class=\"bi bi-star-fill me-2\"></i>Added to favorites!</div><button type=\"button\" class=\"btn-close btn-close-white me-2 m-auto\" data-bs-dismiss=\"toast\" aria-label=\"Close\"></button></div></div></div><div class=\"row g-4\"><div class=\"col-12 col-lg-8\"><div class=\"mb-4\"><div class=\"input-group\"><span class=\"input-group-text\"><i class=\"bi bi-search\"></i></span> <input type=\"search\" placeholder=\"Search prompts...\" class=\"form-control\" hx-get=\"/search\" name=\"q\" hx-trigger=\"keyup changed delay:200ms\" hx-target=\"#prompt-list\" hx-get-oob=\"true\" hx-get-oob-swap=\"true\" hx-get-oob-url=\"/\"></div></div><div id=\"prompt-list\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, repo := range repositories {
				templ_7745c5c3_Err = components.PromptList(repos[repo].GetPromptos()).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"col-12 col-lg-4\"><div id=\"prompt-content\" class=\"card mb-4\"><div class=\"card-body\"><p class=\"text-muted\">Select a prompt to view its details</p></div></div><div class=\"card\"><div class=\"card-header d-flex justify-content-between align-items-center\"><h5 class=\"mb-0\">Favorites</h5></div><div class=\"card-body\" id=\"favorites-list\"><p class=\"text-muted mb-0\">No favorites yet</p></div></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = templates.Layout().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
