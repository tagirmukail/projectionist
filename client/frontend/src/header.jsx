import React from 'react';
import { routes } from './router';
import { Link } from 'react-router-dom';

export const Header = () => (
    <header>
        <nav>
            <ul>
                {routes.map((route, i) => (
                    <li key={i}>
                        <Link to={route.url}>{route.name}</Link>
                    </li>
                ))}
            </ul>
        </nav>
    </header>
)