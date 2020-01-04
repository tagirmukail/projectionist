import { FETCH_LOGIN, FETCH_LOGIN_START } from "../types/actionTypes";
import api from '../api/api';

export const fetchLogin = (form) => async (dispatch, getState) => {
    dispatch({ type: FETCH_LOGIN_START, payload: { data: form } });

    api(getState)
        .post(`/v1/api/login`, form)
        .then((resp) => {
            localStorage.setItem("token", resp.data.user.token);
            dispatch({ type: FETCH_LOGIN, payload: resp.data });
        })
        .catch((error) => console.log(error));
};