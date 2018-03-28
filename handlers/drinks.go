package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coffee-shop/models"
	"github.com/coffee-shop/utils"
)

// DrinkSearchResponse defines the response structure for drink search API
type DrinkSearchResponse struct {
	OffsetPrevious int64           `json:"offset_previous,omitempty"`
	OffsetCurrent  int64           `json:"offset_current,omitempty"`
	Hits           []*models.Drink `json:"hits"`
}

// DrinkCreateHandler handles the create request for drink
func DrinkCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, err, 400)
		return
	}

	var drink models.Drink
	if err = json.Unmarshal(body, &drink); err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, err, 400)
		return
	}

	if err = drink.Validate(); err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, err, 400)
		return
	}

	if err = drink.Create(); err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, err, 500)
		return
	}

	utils.SendResponse(w, fmt.Sprintf(`{"id": "%s"}`, drink.ID), 200)
}

// DrinkDeleteHandler handles the delete request for drink
func DrinkDeleteHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	id := queries.Get("id")
	if id == "" {
		log.Println(utils.ErrDeleteMissingID)
		utils.SendErrorResponse(w, utils.ErrDeleteMissingID, 500)
		return
	}

	drink := &models.Drink{ID: id}
	if err := drink.Delete(); err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, err, 500)
		return
	}

	utils.SendResponse(w, `{"status": "success"}`, 200)
}

// DrinkSearchHandler takes a query and sends back the search response
func DrinkSearchHandler(w http.ResponseWriter, r *http.Request) {
	opts, err := models.NewDrinkSearchOptions(r.URL.Query())
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, err, 400)
		return
	}

	drinks, err := (&models.Drink{}).Query(opts)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, err, 500)
		return
	}

	offsetPrev := int64(0)
	if opts.Offset != nil {
		offsetPrev = *opts.Offset
	}

	resp := DrinkSearchResponse{
		OffsetPrevious: offsetPrev,
		OffsetCurrent:  offsetPrev + int64(len(drinks)),
		Hits:           drinks,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		utils.SendErrorResponse(w, err, 500)
		return
	}

	utils.SendResponse(w, string(respBytes), 200)
}
