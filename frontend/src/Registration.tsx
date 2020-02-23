import React from 'react';
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Icon from "@material-ui/core/Icon";
import Container from "@material-ui/core/Container";
import CardHeader from "@material-ui/core/CardHeader";
import CardActions from "@material-ui/core/CardActions";
import Fab from "@material-ui/core/Fab";
import TextField from "@material-ui/core/TextField";
import {KeyboardDatePicker} from "@material-ui/pickers";
import {InputAdornment} from "@material-ui/core";
import {gql} from "apollo-boost";
import {useMutation} from "@apollo/react-hooks";
import {GraphQLError} from "graphql";
import clsx from "clsx";
import {createBrowserHistory} from "history"

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
	save: {
		marginLeft: 'auto'
	}
}));

const history = createBrowserHistory();

const SUBMIT_REGISTRATION = gql`
mutation UserRegister($details: UserRegistration) {
	users {
		register(details: $details) {
			id
		}
	}
}
`;

function Registration() {
	const classes = useStyles();

	const [submitRegistration, submitResult] = useMutation(SUBMIT_REGISTRATION, {
		onCompleted: (data) => {
			if(data?.users?.register?.id) {
				history.push('/')
			}
		}
	});

	const [dobValue, setDobValue] = React.useState<Date | null>(new Date());
	const [nameValue, setNameValue] = React.useState<String>('');
	const [displayNameValue, setDisplayNameValue] = React.useState<String>('');

	const saveIcon = clsx({
		'fas': true,
		'fa-ellipsis-h': submitResult.loading,
		'fa-save': !submitResult.loading
	});

	function hasNameError(gqlErrors: ReadonlyArray<GraphQLError> | undefined): string {
		if (!gqlErrors) {
			return '';
		}

		return gqlErrors.find(e => e.message.startsWith('nameError:'))?.message || '';
	}

	function hasDisplayNameError(gqlErrors: ReadonlyArray<GraphQLError> | undefined): string {
		if (!gqlErrors) {
			return '';
		}

		return gqlErrors.find(e => e.message.startsWith('displayNameError:'))?.message || '';
	}

	function hasDobError(gqlErrors: ReadonlyArray<GraphQLError> | undefined): string {
		if (!gqlErrors) {
			return '';
		}

		return gqlErrors.find(e => e.message.startsWith('dobError:'))?.message || '';
	}

	return (
		<div className={classes.root}>
			<Container className={classes.container}>
				<Card className={classes.card} elevation={4}>
					<CardHeader className={classes.cardHeader} title="Register"/>
					<CardContent>
						<form onSubmit={async (e) => {
							e.preventDefault();
							await submitRegistration({
								variables: {
									details: {
										name: nameValue,
										displayName: displayNameValue,
										dob: dobValue
									}
								}
							});
						}}>
							<TextField
								error={!submitResult.loading && hasNameError(submitResult.error?.graphQLErrors).length > 0}
								helperText={!submitResult.loading ? hasNameError(submitResult.error?.graphQLErrors) : ''}
								onChange={(e) => setNameValue(e.target.value)} margin="normal" fullWidth required
								id="registration-name" label="Name" InputProps={{
								startAdornment: (
									<InputAdornment position="start"><Icon className="fas fa-user"/></InputAdornment>
								)
							}}/>
							<TextField
								error={!submitResult.loading && hasDisplayNameError(submitResult.error?.graphQLErrors).length > 0}
								helperText={!submitResult.loading ? hasDisplayNameError(submitResult.error?.graphQLErrors) : ''}
								onChange={(e) => setDisplayNameValue(e.target.value)} margin="normal" fullWidth
								id="registration-displayname" label="Display Name" InputProps={{
								startAdornment: (
									<InputAdornment position="start"><Icon
										className="fas fa-user-ninja"/></InputAdornment>
								)
							}}/>
							<KeyboardDatePicker
								fullWidth
								error={!submitResult.loading && hasDobError(submitResult.error?.graphQLErrors).length > 0}
								helperText={!submitResult.loading ? hasDobError(submitResult.error?.graphQLErrors) : ''}
								required
								disableToolbar
								variant="inline"
								format="yyyy-MM-dd"
								margin="normal"
								id="registration-dob"
								label="Date of Birth"
								value={dobValue}
								onChange={(d) => setDobValue(d)}
								KeyboardButtonProps={{
									'aria-label': 'change date',
								}}
							/>
						</form>
					</CardContent>
					<CardActions>
						<Fab color="primary" className={classes.save} aria-label="Save">
							<Icon className={saveIcon}/>
						</Fab>
					</CardActions>
				</Card>
			</Container>
		</div>
	);
}

export default Registration;
