import React from 'react';
import './App.css';
import {
	BrowserRouter as Router,
	Switch,
	Route,
	Redirect,
	Link
} from "react-router-dom";
import Login from "./Login";
import Registration from "./Registration"
import DateFnsUtils from "@date-io/date-fns";
import { MuiPickersUtilsProvider } from '@material-ui/pickers';

// A wrapper for <Route> that redirects to the login
// screen if you're not yet authenticated.
// @ts-ignore
function PrivateRoute({children, ...rest}) {
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

function App() {
	return (
		<MuiPickersUtilsProvider utils={DateFnsUtils}>
			<Router>
				<Switch>
					<Route path="/login">
						<Login/>
					</Route>
					<Route path="/register">
						<Registration/>
					</Route>
					<PrivateRoute path="*">

					</PrivateRoute>
				</Switch>
			</Router>
		</MuiPickersUtilsProvider>
	);
}

export default App;
