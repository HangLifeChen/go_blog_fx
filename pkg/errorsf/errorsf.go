package errorsf

type ErrInter interface {
	Error() string                                 // implements error interface
	Wrap(err error) ErrInter                       // wrap the real error
	GetMessageId() string                          // get error message id
	Params(params map[string]interface{}) ErrInter // replace the parameter {{ .errMsg }} in the internationalization code template
	GetMessageParams() map[string]interface{}      // get the parameter replacement
}

type Errors struct {
	Tips          string                 // not used, only for tip
	MessageId     string                 // error message id
	MessageParams map[string]interface{} // parameter replacement
	Err           error                  // real error
}

func NewErrors(tips, messageId string) ErrInter {
	return &Errors{Tips: tips, MessageId: messageId}
}

func (o *Errors) Error() string {
	if o.Err != nil {
		return o.Err.Error()
	}
	return ""
}

func (o *Errors) Wrap(err error) ErrInter {
	o.Err = err
	return o
}

func (o *Errors) GetMessageId() string {
	return o.MessageId
}

func (o *Errors) Params(params map[string]interface{}) ErrInter {
	o.MessageParams = params
	return o
}

func (o *Errors) GetMessageParams() map[string]interface{} {
	return o.MessageParams
}
