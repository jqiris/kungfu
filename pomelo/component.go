package pomelo

import (
	"errors"
	"reflect"
)

type (
	//Handler represents a message.Message's handler's meta information.
	//Handler represents a message.Message's handler's meta information.
	Handler struct {
		Receiver reflect.Value  // receiver of method
		Method   reflect.Method // method stub
		Type     reflect.Type   // low-level type of method
		IsRawArg bool           // whether the data need to serialize
	}

	// Service implements a specific service, some of it's methods will be
	// called when the correspond events is occurred.
	Service struct {
		Name     string              // name of service
		Type     reflect.Type        // type of the receiver
		Receiver reflect.Value       // receiver of methods for the service
		Handlers map[string]*Handler // registered methods
		Options  options             // options
	}
	// Component is the interface that represent a component.
	Component interface {
		Init()
		AfterInit()
		BeforeShutdown()
		Shutdown()
	}
)

func NewService(comp Component, opts []Option) *Service {
	s := &Service{
		Type:     reflect.TypeOf(comp),
		Receiver: reflect.ValueOf(comp),
	}

	// apply options
	for i := range opts {
		opt := opts[i]
		opt(&s.Options)
	}
	if name := s.Options.name; name != "" {
		s.Name = name
	} else {
		s.Name = reflect.Indirect(s.Receiver).Type().Name()
	}

	return s
}

// suitableMethods returns suitable methods of typ
func (s *Service) suitableHandlerMethods(typ reflect.Type) map[string]*Handler {
	methods := make(map[string]*Handler)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mt := method.Type
		mn := method.Name
		if isHandlerMethod(method) {
			raw := false
			if mt.In(2) == typeOfBytes {
				raw = true
			}
			// rewrite handler name
			if s.Options.nameFunc != nil {
				mn = s.Options.nameFunc(mn)
			}
			methods[mn] = &Handler{Method: method, Type: mt.In(2), IsRawArg: raw}
		}
	}
	return methods
}

// ExtractHandler extract the set of methods from the
// receiver value which satisfy the following conditions:
// - exported method of exported type
// - two arguments, both of exported type
// - the first argument is *session.Session
// - the second argument is []byte or a pointer
func (s *Service) ExtractHandler() error {
	typeName := reflect.Indirect(s.Receiver).Type().Name()
	if typeName == "" {
		return errors.New("no service name for type " + s.Type.String())
	}
	if !isExported(typeName) {
		return errors.New("type " + typeName + " is not exported")
	}

	// Install the methods
	s.Handlers = s.suitableHandlerMethods(s.Type)

	if len(s.Handlers) == 0 {
		str := ""
		// To help the user, see if a pointer receiver would work.
		method := s.suitableHandlerMethods(reflect.PtrTo(s.Type))
		if len(method) != 0 {
			str = "type " + s.Name + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "type " + s.Name + " has no exported methods of suitable type"
		}
		return errors.New(str)
	}

	for i := range s.Handlers {
		s.Handlers[i].Receiver = s.Receiver
	}

	return nil
}

var (
	comps = make([]regComp, 0)
)

type regComp struct {
	comp Component
	opts []Option
}

func startupComponents() {
	// component initialize hooks
	for _, c := range comps {
		c.comp.Init()
	}

	// component after initialize hooks
	for _, c := range comps {
		c.comp.AfterInit()
	}

	// register all components
	for _, c := range comps {
		if err := handler.register(c.comp, c.opts); err != nil {
			logger.Println(err.Error())
		}
	}

	handler.DumpServices()
}

func shutdownComponents() {
	// reverse call `BeforeShutdown` hooks
	length := len(comps)
	for i := length - 1; i >= 0; i-- {
		comps[i].comp.BeforeShutdown()
	}

	// reverse call `Shutdown` hooks
	for i := length - 1; i >= 0; i-- {
		comps[i].comp.Shutdown()
	}
}
