package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"preview-week1/entity"
	"strconv"
)

type NewBranchHandler struct {
	*sql.DB
}

type NewBranch struct {
	ID int `json:"branch_id"`
}

func (h *NewBranchHandler) GetAllBranches(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var branches []entity.Branch

	rows, err := h.Query("SELECT branch_id, name, location FROM branches")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while fetching branches from database",
		})
		log.Println("Error while fetching branches from database:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var branch entity.Branch

		err = rows.Scan(&branch.ID, &branch.Name, &branch.Location)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Message: "Error while scanning branches from database",
			})
			log.Println("Error while scanning branches from database:", err)
			return
		}

		branches = append(branches, branch)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Branches retrieved successfully",
		Data:    branches,
	})
}

func (h *NewBranchHandler) GetBranchById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var branch entity.Branch
	paramsId := p.ByName("id")

	rows, err := h.Query("SELECT branch_id, name, location FROM branches WHERE branch_id=?", paramsId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while fetching branch from database",
		})
		log.Println("Error while fetching branch from database:", err)
		return
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true

		err = rows.Scan(&branch.ID, &branch.Name, &branch.Location)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.Message{
				Status:  "failed",
				Code:    http.StatusInternalServerError,
				Message: "Error while scanning branch from database",
			})
			log.Println("Error while scanning branch from database:", err)
			return
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusNotFound,
			Message: "Branch not found",
		})
		log.Println("Branch not found:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Branch retrieved successfully",
		Data:    branch,
	})
}

func (h *NewBranchHandler) CreateNewBranch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var branch entity.Branch

	err := json.NewDecoder(r.Body).Decode(&branch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Error while parsing request body",
		})
		log.Println("Error while parsing request body:", err)
		return
	}

	result, err := h.Exec("INSERT INTO branches (name ,location) VALUES (?, ?)", branch.Name, branch.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while inserting branch",
		})
		log.Println("Error while inserting branch:", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while fetching branch from database",
		})
		log.Println("Error while fetching branch from database:", err)
		return
	}
	branch.ID = int(id)
	newBranch := NewBranch{
		ID: int(id),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusCreated,
		Message: "Branch created successfully",
		Data:    newBranch,
	})
}

func (h *NewBranchHandler) UpdateBranch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var branch entity.Branch
	paramsId := p.ByName("id")

	err := json.NewDecoder(r.Body).Decode(&branch)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Error while parsing request body",
		})
		log.Println("Error while parsing request body:", err)
		return
	}

	branch.ID, err = strconv.Atoi(paramsId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusBadRequest,
			Message: "Error while parsing request body",
		})
		log.Println("Error while parsing request body:", err)
		return
	}

	_, err = h.Exec(`UPDATE branches SET name = ?, location = ? WHERE branch_id = ?`, branch.Name, branch.Location, branch.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while updating branch",
		})
		log.Println("Error while updating branch:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Branch updated successfully",
	})
}

func (h *NewBranchHandler) DeleteBranch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	paramsId := p.ByName("id")

	_, err := h.Exec(`DELETE FROM branches WHERE branch_id = ?`, paramsId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.Message{
			Status:  "failed",
			Code:    http.StatusInternalServerError,
			Message: "Error while deleting branch",
		})
		log.Println("Error while deleting branch:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.Message{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Branch deleted successfully",
	})
}
