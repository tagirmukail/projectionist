import { FETCH_USERS, FETCH_USERS_START } from '../types/actionTypes'

const initialState = {
    wait: false,
    result: {
        data: [],
    }
};

const users = (state = initialState, action) => {
    switch (action.type) {
        case FETCH_USERS_START:
            return { wait: true, result: { data: [] } };
        case FETCH_USERS:
            return {...state, wait: false, result: action.payload };
        default:
            return state;
    }
};

export default users;