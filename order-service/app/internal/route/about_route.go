package router

import (
	aboutRepo "order-service/app/domain/usercases/about/repo"
	aboutHandler "order-service/app/transport/about/handler"
)

func setupAboutRoute(v1 fiber.Router) {
	groupAboutHandler := aboutHandler.AboutHandler{
		AboutDomain: aboutRepo.NewAboutRepo(),
	}

	groupAbout := v1.Group("/abouts")
	{
		GET(groupAbout, "/:sectionId", groupAboutHandler.GetAboutBySectionId)
	}
}
