import {
	Card,
	CardContent,
	CardHeader,
	Icon,
	List,
	ListItem,
	ListItemIcon,
	ListItemText,
	makeStyles
} from "@material-ui/core";
import firebase from "firebase";

const useStyles = makeStyles({
	root: {
		minWidth: 275,
		maxWidth: 275,
	}
})

function SignIn() {
	const classes = useStyles();

	const handleGoogleAuth = () => {
		const provider = new firebase.auth.GoogleAuthProvider();
		firebase.auth().signInWithPopup(provider)
			.then((result) => {
				result.user?.getIdToken(false).then((token) => {
					console.log(token);
				})
			}).catch((err) => {
				console.log(err);
		})
	}

	return (
		<Card className={classes.root}>
			<CardHeader title="Sign In"/>
			<CardContent>
				<List>
					<ListItem button onClick={handleGoogleAuth}>
						<ListItemIcon>
							<Icon className="fab fa-google"/>
						</ListItemIcon>
						<ListItemText>Google</ListItemText>
					</ListItem>
				</List>
			</CardContent>
		</Card>
	);
}

export default SignIn;