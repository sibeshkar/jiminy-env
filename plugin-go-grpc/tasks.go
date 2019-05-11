package main

var taskList = map[string]string{
	"BisectAngle":         "http://localhost:3000/miniwob/bisect-angle.html",
	"BookFlight":          "http://localhost:3000/miniwob/book-flight.html",
	"ChaseCircle":         "http://localhost:3000/miniwob/chase-circle.html",
	"ChooseDate":          "http://localhost:3000/miniwob/choose-date.html",
	"ChooseList":          "http://localhost:3000/miniwob/choose-list.html",
	"CircleCenter":        "http://localhost:3000/miniwob/circle-center.html",
	"ClickButton":         "http://localhost:3000/miniwob/click-button.html",
	"ClickButtonSequence": "http://localhost:3000/miniwob/click-button-sequence.html",
	"ClickCheckboxes":     "http://localhost:3000/miniwob/click-checkboxes.html",
	"ClickCollapsible":    "http://localhost:3000/miniwob/click-collapsible.html",
	"ClickCollapsible2":   "http://localhost:3000/miniwob/click-collapsible-2.html",
	"ClickColor":          "http://localhost:3000/miniwob/click-color.html",
	"ClickDialog":         "http://localhost:3000/miniwob/click-dialog.html",
	"ClickDialog2":        "http://localhost:3000/miniwob/click-dialog-2.html",
	"ClickLink":           "http://localhost:3000/miniwob/click-link.html",
	"ClickMenu":           "http://localhost:3000/miniwob/click-menu.html",
	"ClickMenu2":          "http://localhost:3000/miniwob/click-menu-2.html",
	"ClickOption":         "http://localhost:3000/miniwob/click-option.html",
	"ClickPie":            "http://localhost:3000/miniwob/click-pie.html",
	"ClickScrollList":     "http://localhost:3000/miniwob/click-scroll-list.html",
	"ClickShades":         "http://localhost:3000/miniwob/click-shades.html",
	"ClickShape":          "http://localhost:3000/miniwob/click-shape.html",
	"ClickTab":            "http://localhost:3000/miniwob/click-tab.html",
	"ClickTab2":           "http://localhost:3000/miniwob/click-tab-2.html",
	"ClickTest":           "http://localhost:3000/miniwob/click-test.html",
	"ClickTest2":          "http://localhost:3000/miniwob/click-test-2.html",
	"ClickWidget":         "http://localhost:3000/miniwob/click-widget.html",
	"CopyPaste":           "http://localhost:3000/miniwob/copy-paste.html",
	"CopyPaste2":          "http://localhost:3000/miniwob/copy-paste-2.html",
	"CountShape":          "http://localhost:3000/miniwob/count-shape.html",
	"CountSides":          "http://localhost:3000/miniwob/count-sides.html",
	"DragBox":             "http://localhost:3000/miniwob/drag-box.html",
	"DragCube":            "http://localhost:3000/miniwob/drag-cube.html",
	"DragItem":            "http://localhost:3000/miniwob/drag-item.html",
	"DragItems":           "http://localhost:3000/miniwob/drag-items.html",
	"DragItemsGrid":       "http://localhost:3000/miniwob/drag-items-grid.html",
	"DragShapes":          "http://localhost:3000/miniwob/drag-shapes.html",
	"DragSortNumbers":     "http://localhost:3000/miniwob/drag-sort-numbers.html",
	"EmailInbox":          "http://localhost:3000/miniwob/email-inbox.html",
	"EnterDate":           "http://localhost:3000/miniwob/enter-date.html",
	"EnterPassword":       "http://localhost:3000/miniwob/enter-password.html",
	"EnterText":           "http://localhost:3000/miniwob/enter-text.html",
	"EnterText2":          "http://localhost:3000/miniwob/enter-text-2.html",
	"EnterTextDynamic":    "http://localhost:3000/miniwob/enter-text-dynamic.html",
	"EnterTime":           "http://localhost:3000/miniwob/enter-time.html",
	"FindMidpoint":        "http://localhost:3000/miniwob/find-midpoint.html",
	"FindWord":            "http://localhost:3000/miniwob/find-word.html",
	"FocusText":           "http://localhost:3000/miniwob/focus-text.html",
	"FocusText2":          "http://localhost:3000/miniwob/focus-text-2.html",
	"GridCoordinate":      "http://localhost:3000/miniwob/grid-coordinate.html",
	"GuessNumber":         "http://localhost:3000/miniwob/guess-number.html",
	"HighlightText":       "http://localhost:3000/miniwob/highlight-text.html",
	"HighlightText2":      "http://localhost:3000/miniwob/highlight-text-2.html",
	"IdentifyShape":       "http://localhost:3000/miniwob/identify-shape.html",
	"LoginUser":           "http://localhost:3000/miniwob/login-user.html",
	"MovingItems":         "http://localhost:3000/miniwob/moving-items.html",
	"NavigateTree":        "http://localhost:3000/miniwob/navigate-tree.html",
	"NumberCheckboxes":    "http://localhost:3000/miniwob/number-checkboxes.html",
	"ReadTable":           "http://localhost:3000/miniwob/read-table.html",
	"ReadTable2":          "http://localhost:3000/miniwob/read-table-2.html",
	"ResizeTextarea":      "http://localhost:3000/miniwob/resize-textarea.html",
	"RightAngle":          "http://localhost:3000/miniwob/right-angle.html",
	"ScrollText":          "http://localhost:3000/miniwob/scroll-text.html",
	"ScrollText2":         "http://localhost:3000/miniwob/scroll-text-2.html",
	"SearchEngine":        "http://localhost:3000/miniwob/search-engine.html",
	"SimonSays":           "http://localhost:3000/miniwob/simon-says.html",
	"SimpleAlgebra":       "http://localhost:3000/miniwob/simple-algebra.html",
	"SimpleArithmetic":    "http://localhost:3000/miniwob/simple-arithmetic.html",
	"SocialMedia":         "http://localhost:3000/miniwob/social-media.html",
	"Terminal":            "http://localhost:3000/miniwob/terminal.html",
	"TextEditor":          "http://localhost:3000/miniwob/text-editor.html",
	"TextTransform":       "http://localhost:3000/miniwob/text-transform.html",
	"TicTacToe":           "http://localhost:3000/miniwob/tic-tac-toe.html",
	"UseAutocomplete":     "http://localhost:3000/miniwob/use-autocomplete.html",
	"UseColorwheel":       "http://localhost:3000/miniwob/use-colorwheel.html",
	"UseColorwheel2":      "http://localhost:3000/miniwob/use-colorwheel-2.html",
	"UseSlider":           "http://localhost:3000/miniwob/use-slider.html",
	"UseSlider2":          "http://localhost:3000/miniwob/use-slider-2.html",
	"UseSpinner":          "http://localhost:3000/miniwob/use-spinner.html",
	"VisualAddition":      "http://localhost:3000/miniwob/visual-addition.html",
}