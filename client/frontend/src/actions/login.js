import { FETCH_LOGIN, FETCH_LOGIN_START } from "../types/actionTypes";
import api from '../api/api';

export const fetchLogin = (form) => async (dispatch, getState) => {
    dispatch({ type: FETCH_LOGIN_START, payload: { data: form } });

    api(getState)
        .post(`/v1/api/login`, form)
        .then((resp) => {
            dispatch({ type: FETCH_LOGIN, payload: resp.data })
            return;
        })
        .catch((error) => console.log(error));
}