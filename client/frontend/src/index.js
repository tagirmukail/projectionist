import React from 'react';
import ReactDOM from 'react-dom';
import { createStore, applyMiddleware, combineReducers } from 'redux';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';
import { BrowserRouter } from 'react-router-dom';
import MainRouter from './router';
import './index.css';
import reducers from './reducers/reducers'
import { Header } from './header';

const store = createStore(combineReducers(reducers), applyMiddleware(thunk));

ReactDOM.render(
    <Provider store={store}>
        <BrowserRouter>
            <div>
                <Header />
                <MainRouter />
            </div>
        </BrowserRouter>
    </Provider>, document.getElementById('root')
);