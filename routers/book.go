package routers

import (
	"MetaLib/models"
	"MetaLib/templmanager"
	"MetaLib/utils"
	"errors"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func bookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error(err)
	}

	book, bookErr := models.GetBookById(bookId)

	var userRatingNumber int

	session := utils.GetSession(r)
	user, err := models.GetUser(session)
	if err != nil || user.Status != models.Logged {
		userRatingNumber = 0
	} else {
		var userRating models.Rating
		userRatingNotFound := utils.DB.Where("user_id = ? and book_id = ?", user.Id, bookId).First(&userRating).RecordNotFound()

		if userRatingNotFound {
			userRatingNumber = 0
		} else {
			userRatingNumber = userRating.Rating
		}
	}

	ratingStars := make([]int, 10)
	for i := range ratingStars {
		ratingStars[i] = 10 - i
	}

	err = templmanager.RenderTemplate(w, r, "book.html", struct {
		Book        *models.Book
		BookErr     error
		UserRating  int
		RatingStars []int
	}{Book: book, BookErr: bookErr, UserRating: userRatingNumber, RatingStars: ratingStars})
	if err != nil {
		log.Fatal(err)
	}
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	perPage := 4
	displayPages := 7

	var count int
	utils.DB.Model(&models.Book{}).Count(&count)

	pages := count/perPage + 1
	if count%perPage == 0 {
		pages--
	}

	pageString := r.URL.Query().Get("page")
	if len(pageString) == 0 {
		pageString = "1"
	}

	var booksErr error
	page, err := strconv.Atoi(pageString)
	if err != nil {
		booksErr = errors.New("bad page")
	}

	startFrom := page - displayPages/2
	if startFrom <= 0 {
		startFrom = 1
	}

	endTo := page + displayPages/2
	if endTo > pages {
		endTo = pages
	}

	pagesSqn := make([]int, endTo-startFrom+1)

	for i := startFrom; i <= endTo; i++ {
		ind := i - startFrom
		pagesSqn[ind] = i
	}

	var books []models.Book
	offset := (page - 1) * perPage
	notFind := utils.DB.Limit(perPage).Offset(offset).Set("gorm:auto_preload", true).Find(&books).RecordNotFound()
	if notFind {
		booksErr = errors.New("books not found")
	}

	var nextPage, previousPage int
	if page+1 > pages {
		nextPage = page
	} else {
		nextPage = page + 1
	}
	if page-1 > 0 {
		previousPage = page - 1
	} else {
		previousPage = page
	}

	err = templmanager.RenderTemplate(w, r, "books.html", struct {
		Books        []models.Book
		BooksErr     error
		Page         int
		Pages        []int
		NextPage     int
		PreviousPage int
		PageCount    int
	}{Books: books, BooksErr: booksErr, Page: page, Pages: pagesSqn, NextPage: nextPage, PreviousPage: previousPage, PageCount: pages})
	if err != nil {
		log.Fatal(err)
	}
}

func starHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error(err)
		return
	}

	ratingStr := r.PostFormValue("rating")
	if len(ratingStr) < 0 {
		return
	}
	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		log.Error(err)
		return
	}

	bookIdSrt := r.PostFormValue("book-id")
	if len(bookIdSrt) < 0 {
		return
	}
	bookId, err := strconv.Atoi(bookIdSrt)
	if err != nil {
		log.Error(err)
		return
	}

	user, err := models.GetUserR(r)
	if err != nil {
		log.Error(err)
		return
	}

	err = utils.DB.Where(models.Rating{UserId: user.Id, BookId: bookId}).Assign(models.Rating{Rating: rating}).FirstOrCreate(&models.Rating{}).Error
	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusOK)
}
