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

function Registration() {
	const classes = useStyles();

	const [selectedDate, setSelectedDate] = React.useState<Date | null>(
		new Date(),
	);

	const handleDateChange = (date: Date | null) => {
		setSelectedDate(date);
	};

	return (
		<div className={classes.root}>
			<Container className={classes.container}>
				<Card className={classes.card} elevation={4}>
					<CardHeader className={classes.cardHeader} title="Register"/>
					<CardContent>
						<form>
							<TextField margin="normal" fullWidth required id="registration-name" label="Name" InputProps={{
								startAdornment: (
									<InputAdornment position="start"><Icon className="fas fa-user" /></InputAdornment>
								)
							}}/>
							<TextField margin="normal" fullWidth id="registration-displayname" label="Display Name" InputProps={{
								startAdornment: (
									<InputAdornment position="start"><Icon className="fas fa-user-ninja"/></InputAdornment>
								)
							}}/>
							<KeyboardDatePicker
								fullWidth
								required
								disableToolbar
								variant="inline"
								format="yyyy-MM-dd"
								margin="normal"
								id="registration-dob"
								label="Date of Birth"
								value={selectedDate}
								onChange={handleDateChange}
								KeyboardButtonProps={{
									'aria-label': 'change date',
								}}
							/>
						</form>
					</CardContent>
					<CardActions>
						<Fab color="primary" className={classes.save} aria-label="Save">
							<Icon className="fas fa-save"/>
						</Fab>
					</CardActions>
				</Card>
			</Container>
		</div>
	);
}

export default Registration;
