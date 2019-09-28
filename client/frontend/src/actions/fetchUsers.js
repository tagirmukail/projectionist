import { FETCH_USERS, FETCH_USERS_START } from "../types/actionTypes";
import api from '../api/api'
// import { handleError } from "./errors";

export const fetchUsers = (page, count) => async (dispatch, getState) => {
    dispatch({ type: FETCH_USERS_START });

    api(getState)
        .get(`/v1/api/user?page=${page}&count=${count}`)
        .then((resp) => {
            console.log(resp);
            dispatch({type: FETCH_USERS, payload: resp.data})
        })
        .catch((error) => {
            console.log(error);
        })
}