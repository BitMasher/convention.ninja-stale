import React from 'react';
import {
	BrowserRouter as Router,
	Switch,
	Route,
} from "react-router-dom";
import InventoryHome from "./InventoryHome";
import InventoryImport from "./InventoryImport";

function InventoryApp() {
	return (
		<Router>
			<Switch>
				<Route path="/inventory/import">
					<InventoryImport/>
				</Route>
				<Route path="/inventory">
					<InventoryHome/>
				</Route>
			</Switch>
		</Router>
	)
}

export default InventoryApp;
