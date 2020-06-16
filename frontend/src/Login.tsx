import React from 'react';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import {List, ListItemProps} from "@material-ui/core";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import Icon from "@material-ui/core/Icon";
import Container from "@material-ui/core/Container";
import CardHeader from "@material-ui/core/CardHeader";

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

function ListItemLink(props: ListItemProps<'a', { button?: true }>) {
	return <ListItem button component="a" {...props} />;
}

function Login() {
	const classes = useStyles();

	return (
		<div className={classes.root}>
			<Container className={classes.container}>
				<Card className={classes.card} elevation={4}>
					<CardHeader className={classes.cardHeader} title="Login"/>
					<CardContent>
						<List component="nav" aria-label="google facebook">
							<ListItemLink href="/auth/google">
								<ListItemIcon><Icon className="fab fa-google"/></ListItemIcon>
								<ListItemText>Sign in with Google</ListItemText>
							</ListItemLink>
							<ListItemLink href="/auth/facebook">
								<ListItemIcon><Icon className="fab fa-facebook"/></ListItemIcon>
								<ListItemText>Sign in with Facebook</ListItemText>
							</ListItemLink>
						</List>
					</CardContent>
				</Card>
			</Container>
		</div>
	);
}

export default Login;
