import React from 'react';
import './App.css';
import {
	BrowserRouter as Router,
	Switch,
	Route,
	Redirect
} from "react-router-dom";
import Login from "./Login";
import Registration from "./Registration"
import DateFnsUtils from "@date-io/date-fns";
import {MuiPickersUtilsProvider} from '@material-ui/pickers';
import {ApolloProvider} from '@apollo/react-hooks';
import Cookie from 'js-cookie';

import ApolloClient from 'apollo-boost';
import Cookies from "js-cookie";

const client = new ApolloClient({
	uri: '/graphql',
});

// A wrapper for <Route> that redirects to the login
// screen if you're not yet authenticated.
// @ts-ignore
function PrivateRoute({children, ...rest}) {
	return (
		<Route
			{...rest}
			render={({location}) => {
				const redirect = <Redirect to={{pathname: "/login"}}/>;
				const regToken = Cookies.get('token');
				let payload: { name: string, aud: string, exp: number } | null = null;
				if (regToken) {
					const [, payloadb64,] = regToken.split('.');
					payload = JSON.parse(window.atob(payloadb64));
				}
				if (!payload) {
					return redirect;
				}
				let exp = new Date(payload.exp * 1000);
				if (payload.aud !== 'api' || exp < (new Date())) {
					return redirect;
				}
				return children;
			}
			}
		/>
	);
}

// @ts-ignore
function RegisterRestrictedRoute({children, ...rest}) {
	return (
		<Route
			{...rest}
			render={({location}) =>
				Cookie.get('token') ? (
					children
				) : (
					<Redirect
						to={{
							pathname: "/login"
						}}
					/>
				)
			}
		/>
	);
}


function App() {
	return (
		<ApolloProvider client={client}>
			<MuiPickersUtilsProvider utils={DateFnsUtils}>
				<Router>
					<Switch>
						<Route path="/login">
							<Login/>
						</Route>
						<RegisterRestrictedRoute path="/register">
							<Registration/>
						</RegisterRestrictedRoute>
						<PrivateRoute path="*">
							congrats you've been authorized
						</PrivateRoute>
					</Switch>
				</Router>
			</MuiPickersUtilsProvider>
		</ApolloProvider>
	);
}

export default App;
