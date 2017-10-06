package manipulador

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jeffprestes/cursodego/banco_sql/model"
	"github.com/jeffprestes/cursodego/banco_sql/repo"
)

//Local é o manipulador da requisição de rota /local/
func Local(w http.ResponseWriter, r *http.Request) {
	local := model.Local{}
	param := r.URL.Path[7:]
	codigoTelefone, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "Não foi enviado um numero válido. Verifique.", http.StatusBadRequest)
		fmt.Println("[local] erro ao converter o numero enviado: ", err.Error())
		return
	}
	db, err := repo.GetDB()
	if err != nil {
		fmt.Println("[Local] Erro ao conectar ao Oracle: ", err.Error())
		http.Error(w, "Erro interno no servidor. Verifique.", http.StatusInternalServerError)
		return
	}
	sql := "select country, city, telcode from place where telcode = :codigoTel"
	err = db.Get(&local, sql, codigoTelefone)
	if err != nil {
		http.Error(w, "Não foi possível pesquisar esse numero.", http.StatusInternalServerError)
		fmt.Println("[local] nao foi possível executar a query: ", sql, " Erro: ", err.Error())
		return
	}
	if err := ModeloLocal.ExecuteTemplate(w, "local.html", local); err != nil {
		http.Error(w, "Houve um erro na renderização da página.", http.StatusInternalServerError)
		fmt.Println("[local] Erro na execucao do modelo: ", err.Error())
	}
	sql = "insert into logquery (daterequest) values (TO_DATE(:param, 'yyyy-mm-dd hh24:mi:ss'))"
	resultado, err := db.Exec(sql, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("[local] Erro na inclusao do log: ", sql, " - ", err.Error())
		return
	}
	linhasAfetadas, err := resultado.RowsAffected()
	if err != nil {
		fmt.Println("[local] Erro ao pegar o numero de linhas afetadas pelo comando: ", sql, " - ", err.Error())
	}
	fmt.Println("Sucesso! ", linhasAfetadas, " linha(s) afetada(s)")
}
