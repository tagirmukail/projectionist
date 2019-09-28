export const handleError = (dispatch, type, data) => (error) => {
    console.log(error);

    if (error.response) {
        const {message, state} = error.response.data;
        dispatch({type, payload: {message, state, data}});
        return;
    }

    dispatch({
        type,
        payload: {
            message: "Request with error",
            state: false,
            data,
        },
    });
};