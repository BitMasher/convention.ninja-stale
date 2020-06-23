import React from 'react';
import {AppBar, Toolbar, Typography} from "@material-ui/core";

function AppPortalTopBar() {
	return (
		<AppBar>
			<Toolbar>
				<Typography variant="h6">
					Convention.Ninja
				</Typography>
			</Toolbar>
		</AppBar>
	);
}

export default AppPortalTopBar;
