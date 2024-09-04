package request

type SelectorUUID struct {
	ID string `uri:"id" binding:"required,uuid"`
}
