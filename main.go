package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ophum/etags-optimistic-concurrency/entities"
	"github.com/ophum/etags-optimistic-concurrency/store"
)

var (
	s *store.Store
)

func main() {
	r := gin.Default()

	s = store.New()

	pets := r.Group("/api/v1/pets")
	pets.GET("", findAll)
	pets.GET("/:id", find)
	pets.POST("", create)
	pets.PUT("/:id", update)
	pets.DELETE("/:id", delete)
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func findAll(ctx *gin.Context) {
	pets, err := s.GetAll()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, pets)
}

func find(ctx *gin.Context) {
	id := ctx.Param("id")
	pet, err := s.Get(id)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	hash := generateHash(pet)
	ctx.Header("ETag", hash)

	ctx.JSON(http.StatusOK, pet)
}

type PetCreateRequest struct {
	Name string `json:"name"`
}

func create(ctx *gin.Context) {
	var req PetCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pet := entities.Pet{
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	pet.GenerateID()

	if err := s.Put(&pet); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, pet)
}

type PetUpdateRequest struct {
	Name string `json:"name"`
}

func update(ctx *gin.Context) {
	var req PetUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := ctx.Param("id")

	pet, err := s.Get(id)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	// clientが持っているhashと現在のentityのhashを比較する
	// 一致していればclientが取得したときから変更されていない
	// 一致していなければclientが取得したときから変更されていることになるためエラーにする
	ifMatch := ctx.GetHeader("If-Match")
	hash := generateHash(pet)
	log.Println("If-Match   ", ifMatch)
	log.Println("entity hash", hash)
	if ifMatch != hash {
		ctx.AbortWithStatus(http.StatusPreconditionFailed)
		return
	}

	pet.Name = req.Name
	pet.UpdatedAt = time.Now()

	if err := s.Put(pet); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, pet)

}

func delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if _, err := s.Get(id); err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	s.Delete(id)
	ctx.Status(http.StatusNoContent)
}

func generateHash(v any) string {
	j, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(j)
	return hex.EncodeToString(hash[:])
}
