package main

import (
	"log"

	"go.uber.org/fx"
)

func main() {
	// t := Title("Lee Jun ki")
	// p := NewPublisher(&t)
	// m := NewMainService(p)
	// m.Run()
	fx.New(
		fx.Provide(NewMainService),
		fx.Provide(
			fx.Annotate(
				NewPublisher,
				fx.As(new(IPublisher)),
				fx.ParamTags(`group:"titles"`),
			),
		),
		fx.Provide(NewPublisher),
		fx.Provide(
			titleComponent("Goodbye"),
		),
		fx.Provide(
			titleComponent("World"),
		),
		fx.Invoke(func(service *MainService) {
			service.Run()
		}),
	).Run()
}

func titleComponent(title string) any {
	return fx.Annotate(
		func() *Title {
			t := Title(title)
			return &t
		},
		fx.ResultTags(`group:"titles"`),
	)
}

type MainService struct {
	publisher IPublisher
}

func NewMainService(publisher IPublisher) *MainService {
	return &MainService{publisher: publisher}
}

func (service *MainService) Run() {
	service.publisher.Publish()
	log.Print("Main program")
}

// Dependency

type IPublisher interface {
	Publish()
}

type Publisher struct {
	titles []*Title
}

func NewPublisher(titles ...*Title) *Publisher {
	return &Publisher{titles: titles}
}

func (publisher *Publisher) Publish() {
	for _, title := range publisher.titles {
		log.Print("publisher: ", *title)
	}

}

type Title string
