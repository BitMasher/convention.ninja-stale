import React, {useEffect} from 'react';
import {
	Card, CardContent,
	Container,
	FormGroup,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow, TextField
} from "@material-ui/core";
import {useLazyQuery} from "@apollo/react-hooks";
import {gql} from "apollo-boost";

const ASSET_SEARCH = gql`
query AssetSearch($term: String) {
	assets {
		search(term: $term) {
			id,
			category,
			model,
			manufacturer,
			location
		}
	}
}
`;

function debouncer(delay: number): Function {
	let timeout: number = 0;
	return (cb: Function) => {
		clearTimeout(timeout);
		timeout = setTimeout(cb, delay);
	}
}

const debounce = debouncer(300);
function InventoryHome() {
	const [doSearch, searchResult] = useLazyQuery(ASSET_SEARCH);

	useEffect(() => {
		console.log(searchResult);
	}, [searchResult]);

	return (
		<Container>
			<Card>
				<CardContent>
					<FormGroup>
						<TextField label="Search" onChange={(evt) => {
							const term = evt.target.value;
							debounce(async () => {
								await doSearch({
									variables: {
										term: term
									}
								});
							});
						}}/>
					</FormGroup>
				</CardContent>
			</Card>
			<TableContainer>
				<Table>
					<TableHead>
						<TableRow>
							<TableCell>Asset ID</TableCell>
							<TableCell>Category</TableCell>
							<TableCell>Model</TableCell>
							<TableCell>Manufacturer</TableCell>
							<TableCell>Location</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						<TableRow>
							<TableCell>1</TableCell>
							<TableCell>Hotspot</TableCell>
							<TableCell>Hotspot</TableCell>
							<TableCell>AT&T</TableCell>
							<TableCell>Storage 1</TableCell>
						</TableRow>
						<TableRow>
							<TableCell>2</TableCell>
							<TableCell>Laptop</TableCell>
							<TableCell>Insperon</TableCell>
							<TableCell>Dell</TableCell>
							<TableCell>Storage 2</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</TableContainer>
		</Container>
	)
}

export default InventoryHome;
