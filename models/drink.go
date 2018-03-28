package models

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/coffee-shop/db"
	"github.com/coffee-shop/utils"
)

// Drink represents the drink model
type Drink struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Price       string     `json:"price"`
	Start       *time.Time `json:"start"`
	End         *time.Time `json:"end"`
	Ingredients []string   `json:"ingredients"`
}

// DrinkSearchOptions declares the search parameters for querying drinks
type DrinkSearchOptions struct {
	Name        *string
	Date        *time.Time
	Offset      *int64
	Limit       *int64
	Ingredients *[]string
}

// NewDrinkSearchOptions creates a new DrinkSearchOptions object from the given query params
func NewDrinkSearchOptions(params url.Values) (DrinkSearchOptions, error) {
	opts := DrinkSearchOptions{}

	if name := params.Get("name"); name != "" {
		opts.Name = &name
	}

	if dateStr := params.Get("date"); dateStr != "" {
		date, err := time.Parse(utils.TimeFormat, dateStr)
		if err != nil {
			return opts, errors.Wrap(err, "Failed parsing date field as a drink search parameter")
		}
		opts.Date = &date
	}

	if offsetStr := params.Get("offset"); offsetStr != "" {
		offset, err := strconv.ParseInt(offsetStr, 0, 64)
		if err != nil {
			return opts, errors.Wrap(err, "Failed parsing offset field as a drink search parameter")
		}

		if offset < 0 {
			return opts, errors.New("Offset must be a non-negative integer")
		}

		opts.Offset = &offset
	}

	if limitStr := params.Get("limit"); limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 0, 64)
		if err != nil {
			return opts, errors.Wrap(err, "Failed parsing limit field as a drink search parameter")
		}

		if limit < 0 {
			return opts, errors.New("Limit must be a non-negative integer")
		}

		opts.Limit = &limit
	}

	if ingrStr := params.Get("ingredients"); ingrStr != "" {
		ingredients := strings.Split(ingrStr, ",")
		opts.Ingredients = &ingredients
	}

	return opts, nil
}

// Validate checks the sanity of the drink object
func (d *Drink) Validate() error {
	err := utils.ValidateMoney(d.Price)
	if err != nil {
		return errors.Wrap(err, "Validation for price failed")
	}

	if d.Start != nil && d.End != nil && d.End.Before(*d.Start) {
		return utils.ErrStartEndDate
	}

	return nil
}

// Create inserts the drink to the DB
func (d *Drink) Create() error {
	tx, err := db.Tx()
	if err != nil {
		return errors.Wrap(err, "Failed initializaing transaction")
	}

	defer tx.Rollback()

	id, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "Failed generating UUID for drink")
	}

	d.ID = id.String()

	var start string
	if d.Start != nil {
		start = d.Start.Format(time.RFC3339)
	}

	var end string
	if d.End != nil {
		end = d.End.Format(time.RFC3339)
	}

	_, err = tx.Exec(`
		INSERT INTO drinks (
			id,
			name,
			price,
			start,
			"end"
		) VALUES ($1, $2, $3, $4, $5)`,
		d.ID,
		d.Name,
		d.Price,
		start,
		end)
	if err != nil {
		return errors.Wrap(err, "Drink insert failed")
	}

	if len(d.Ingredients) > 0 {
		var rows []string
		for _, ingr := range d.Ingredients {
			rows = append(rows, fmt.Sprintf("('%s', '%s')", d.ID, ingr))
		}
		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO drinks_ingredients (
				drink_id,
				ingredient
			) VALUES %s`,
			strings.Join(rows, ", ")))
		if err != nil {
			return errors.Wrap(err, "Drink-Ingredient insert failed")
		}
	}

	return tx.Commit()
}

// Delete deletes the drink row from the DB
func (d *Drink) Delete() error {
	if d.ID == "" {
		return errors.Wrap(utils.ErrDeleteMissingID, "ID missing for deleting drink")
	}

	tx, err := db.Tx()
	if err != nil {
		return errors.Wrap(err, "Failed initializaing transaction")
	}

	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM drinks
		WHERE id = $1`,
		d.ID)
	if err != nil {
		return errors.Wrap(err, "Failed deleting from the drinks table")
	}

	_, err = tx.Exec(`
		DELETE FROM drinks_ingredients
		WHERE drink_id = $1`,
		d.ID)
	if err != nil {
		return errors.Wrap(err, "Failed deleting from the drinks-ingredients table")
	}

	return tx.Commit()
}

// Query searches the drinks with the given search options
func (d *Drink) Query(o DrinkSearchOptions) ([]*Drink, error) {
	drinks := make([]*Drink, 0)

	rows, err := db.DB.Query(dbQuery(o))
	if err != nil {
		return drinks, errors.Wrap(err, "Failed querying drinks")
	}

	for rows.Next() {
		d := &Drink{}
		var ingr string
		err = rows.Scan(
			&d.ID,
			&d.Name,
			&d.Price,
			&d.Start,
			&d.End,
			&ingr)
		if err != nil {
			return drinks, errors.Wrap(err, "Failed reading a drink row")
		}
		d.Ingredients = strings.Split(ingr, ",")

		drinks = append(drinks, d)
	}

	return drinks, nil
}

func dbQuery(o DrinkSearchOptions) string {
	var nameCond string
	if o.Name != nil {
		nameCond = fmt.Sprintf("AND d.name = '%s'", *o.Name)
	}

	var dateCond string
	if o.Date != nil {
		date := (*o.Date).Format(time.RFC3339)
		dateCond = fmt.Sprintf("AND '%s' BETWEEN d.start AND d.end", date)
	}

	var ingrCond string
	if o.Ingredients != nil && len(*o.Ingredients) > 0 {
		var quotedIngrs []string
		for _, ingr := range *o.Ingredients {
			quotedIngrs = append(quotedIngrs, fmt.Sprintf("'%s'", ingr))
		}
		ingrCond = fmt.Sprintf(`
			INNER JOIN drinks_ingredients di2 ON d.id = di2.drink_id WHERE di2.ingredient IN (%s)`,
			strings.Join(quotedIngrs, ","))
	}

	var offsetCond string
	if o.Offset != nil {
		offsetCond = fmt.Sprintf("OFFSET %d", *o.Offset)
	}

	var limitCond string
	if o.Limit != nil {
		limitCond = fmt.Sprintf("LIMIT %d", *o.Limit)
	}

	return fmt.Sprintf(`
		SELECT
			id,
			name,
			price,
			start,
			"end",
			string_agg(DISTINCT(di.ingredient), ',')
		FROM drinks d
		LEFT JOIN drinks_ingredients di ON d.id = di.drink_id
		%s
		%s
		%s
		GROUP BY d.id
		ORDER BY created ASC
		%s
		%s`,
		nameCond,
		dateCond,
		ingrCond,
		offsetCond,
		limitCond)
}
