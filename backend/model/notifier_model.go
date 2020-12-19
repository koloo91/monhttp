package model

type Notifier struct {
	Id      string
	Name    string
	Enabled bool
	Form    []NotificationForm
}

type NotificationForm struct {
	Type            string // the html input type (text, password, email)
	Title           string // include a title for ease of use
	FormControlName string
	Placeholder     string // add a placeholder for the input
	Required        bool   // require this input on the html form
}

type NotifierVo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Form []NotificationFormVo
}

type NotificationFormVo struct {
	Type            string `json:"type"`
	Title           string `json:"title"`
	FormControlName string `json:"formControlName"`
	Placeholder     string `json:"placeholder"`
	Required        bool   `json:"required"`
}

type Notify interface {
	GetId() string
	SendSuccess(Service) error
	SendFailure(Service, Failure) error
	IsEnabled() bool
	GetForms() []NotificationForm
	GetName() string
}

func MapNotifierToVo(n Notify) NotifierVo {
	forms := make([]NotificationFormVo, 0, len(n.GetForms()))
	for _, form := range n.GetForms() {
		forms = append(forms, mapNotificationFormToVo(form))
	}

	return NotifierVo{
		Id:   n.GetId(),
		Name: n.GetName(),
		Form: forms,
	}
}

func mapNotificationFormToVo(n NotificationForm) NotificationFormVo {
	return NotificationFormVo{
		Type:            n.Type,
		Title:           n.Title,
		FormControlName: n.FormControlName,
		Placeholder:     n.Placeholder,
		Required:        n.Required,
	}
}

func MapNotifiersToVos(notifies []Notify) []NotifierVo {
	result := make([]NotifierVo, 0, len(notifies))

	for _, notify := range notifies {
		result = append(result, MapNotifierToVo(notify))
	}

	return result
}
