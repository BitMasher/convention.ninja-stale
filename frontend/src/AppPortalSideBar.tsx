import React from 'react';
import {Drawer, List, ListItem, ListItemIcon, ListItemText} from "@material-ui/core";
import Icon from "@material-ui/core/Icon";
import {makeStyles} from "@material-ui/core/styles";
import {BrowserRouter} from "react-router-dom";

const useStyles = makeStyles(theme => ({
	drawer: {
		width: 240,
		[theme.breakpoints.up('lg')]: {
			marginTop: 64,
			height: 'calc(100% - 64px)'
		}
	},
	nested: {
		paddingLeft: theme.spacing(4)
	},
	divider: {
		margin: theme.spacing(2, 0)
	},
	nav: {
		marginBottom: theme.spacing(2)
	}
}));

function ListItemLink(props: any) {
	return <ListItem button component="a" {...props} />;
}

function AppPortalSideBar() {
	const classes = useStyles();

	return (
		<Drawer classes={{paper: classes.drawer}} anchor="left" variant="persistent" open={true}>
			<div>
				<List>
					<ListItemLink href={"/inventory"}>
						<ListItemIcon>
							<Icon className="fas fa-boxes" fontSize="small"/>
						</ListItemIcon>
						<ListItemText>
							Inventory
						</ListItemText>
					</ListItemLink>
					<List component="div">
						<ListItemLink href={"/inventory/import"} className={classes.nested}>
							<ListItemIcon>
								<Icon className="fas fa-file-import" fontSize="small"/>
							</ListItemIcon>
							<ListItemText>
								Import
							</ListItemText>
						</ListItemLink>
					</List>
				</List>
			</div>
		</Drawer>
	);
}

export default AppPortalSideBar;
