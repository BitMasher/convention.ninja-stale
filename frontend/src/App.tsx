import {Container} from '@material-ui/core';
import React from 'react';
import './App.css';
import SignIn from "./pages/SignIn";
import {ApolloClient, ApolloLink, ApolloProvider, HttpLink, InMemoryCache} from "@apollo/client";

const httpLink = new HttpLink({uri: '/graphql'});

const authLink = new ApolloLink((operation, forward) => {
	// Retrieve the authorization token from local storage.
	const token = localStorage.getItem('auth_token');
	const idToken = localStorage.getItem('id_token');
	// Use the setContext method to set the HTTP headers.
	operation.setContext({
		headers: {
			'X-Identity': idToken ? idToken : '',
			'Authorization': token ? `Bearer ${token}` : '',
		}
	});

	// Call the next link in the middleware chain.
	return forward(operation);
});

const client = new ApolloClient({
	link: authLink.concat(httpLink),
	cache: new InMemoryCache()
});

function App() {

	return (
		<ApolloProvider client={client}>
			<Container>
				<SignIn></SignIn>
			</Container>
		</ApolloProvider>
	);
}

export default App;
