import React from 'react';
import {createStyles, makeStyles, Theme} from "@material-ui/core/styles";
import {Container} from "@material-ui/core";
import AppPortalSideBar from "./AppPortalSideBar";
import AppPortalTopBar from "./AppPortalTopBar";
import {
	BrowserRouter as Router,
	Switch,
	Route,
	Redirect
} from "react-router-dom";
import InventoryApp from "./InventoryApp";

const useStyles = makeStyles((theme: Theme) => createStyles({
	container: {
		textAlign: 'center',
		height: '100vh',
		minHeight: '100vh',
		display: 'flex',
		alignItems: 'center',
		justifyContent: 'center',
	},
	root: {
		width: '100vw',
		height: '100vh',
		backgroundColor: theme.palette.primary.light
	},
	card: {
		minWidth: 275,
		maxWidth: 300,
		margin: 'auto'
	},
	cardHeader: {
		backgroundColor: theme.palette.secondary.light
	},
}));

function AppPortal() {
	const classes = useStyles();
	return (
		<Container className={classes.container}>
			<AppPortalTopBar/>
			<AppPortalSideBar/>
			<main>
				<Router>
					<Switch>
						<Route path="/inventory">
							<InventoryApp/>
						</Route>
						<Route path="">
							<Redirect to="/inventory"/>
						</Route>
					</Switch>
				</Router>
			</main>
		</Container>
	)
}

export default AppPortal;
