package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cyph3rk/api-go-gin/controllers"
	"github.com/cyph3rk/api-go-gin/database"
	"github.com/cyph3rk/api-go-gin/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var ID int

func SetupDasRotasDeTeste() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	rotas := gin.Default()
	return rotas
}

func TestVerificaStatusCodeSaudacaoParametro(t *testing.T) {
	r := SetupDasRotasDeTeste()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/diego", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)

	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")

	mockDaResposta := `{"API dix:":"E ai diego, tudo beleza?"}`
	respostaBody, _ := io.ReadAll(resposta.Body)

	assert.Equal(t, mockDaResposta, string(respostaBody))
}

func CriaLunoMock() {
	aluno := models.Aluno{Nome: "Aluno Teste", CPF: "12345678900", RG: "123456789"}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)
}

func DeletaAlunoMock() {
	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
}

func TestListandoTodosAlunosHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaLunoMock()
	defer DeletaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
}

func TestBuscaAlunoPorCPF(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaLunoMock()
	defer DeletaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	req, _ := http.NewRequest("GET", "/alunos/cpf/12345678900", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
}

func TestBuscaAlunoPorIdhandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaLunoMock()
	defer DeletaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)
	pathDaBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", pathDaBusca, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)

	var alunoMock models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoMock)

	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")

	assert.Equal(t, "Aluno Teste", alunoMock.Nome)
	assert.Equal(t, "12345678900", alunoMock.CPF)
	assert.Equal(t, "123456789", alunoMock.RG)
}

func TestDeletaAlunohandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaLunoMock()
	r := SetupDasRotasDeTeste()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	pathDaBusca := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", pathDaBusca, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)

	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
}

func TestEditAlunohandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaLunoMock()
	defer DeletaAlunoMock()
	r := SetupDasRotasDeTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)

	aluno := models.Aluno{Nome: "Aluno Teste", CPF: "99945678900", RG: "000456789"}
	valorJson, _ := json.Marshal(aluno)

	pathEdicao := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", pathEdicao, bytes.NewBuffer(valorJson))
	resposta := httptest.NewRecorder()
	r.ServeHTTP(resposta, req)

	var alunoAtualizado models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoAtualizado)

	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")

	assert.Equal(t, aluno.Nome, alunoAtualizado.Nome)
	assert.Equal(t, aluno.CPF, alunoAtualizado.CPF)
	assert.Equal(t, aluno.RG, alunoAtualizado.RG)
}
