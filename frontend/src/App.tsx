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

const client = new ApolloClient({
	uri: '/graphql',
});

// A wrapper for <Route> that redirects to the login
// screen if you're not yet authenticated.
// @ts-ignore
function PrivateRoute({children, ...rest}) {
	// TODO: update to check token validity
	return (
		<Route
			{...rest}
			render={({location}) =>
				false ? (
					children
				) : (
					<Redirect
						to={{
							pathname: "/login",
							state: {from: location}
						}}
					/>
				)
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
