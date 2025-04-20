package services

import "net/http"

func (u *ServiceCallsHandler) DeleteServiceById(rw http.ResponseWriter, req *http.Request) {
	id, ok := u.IdValidator("id", req)
	if !ok {
		http.Error(rw, "Invalid Request", http.StatusBadRequest)
		return
	}

	err := u.Repository.DeleteService(id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)

}
