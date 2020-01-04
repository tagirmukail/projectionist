import { FETCH_LOGIN, FETCH_LOGIN_START } from "../types/actionTypes";

const initialState = {
    wait: false,
    result: {
        data: [],
    }
};

export default function login(state = initialState, action) {
    switch (action.type) {
        case FETCH_LOGIN_START:
            return { wait: true, result: { data: [] } };
        case FETCH_LOGIN:
            return {...state, wait: false, result: {data: action.payload} };
        default:
            return state;
    }
};
