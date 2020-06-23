import React, {ChangeEvent, useState} from 'react';
import {
	Button,
	Card,
	CardContent,
	Container,
	FormControlLabel,
	Switch,
} from "@material-ui/core";

enum ImportField {
	BARCODE,
	CATEGORY,
	MODEL,
	MANUFACTURER,
	UNUSED
}

enum ImportStatus {
	READING,
	UNMAPPED,
	CONFLICT,
	FORMAT,
	PROBLEM,
	IMPORTING,
	COMPLETE
}

interface FieldMap {
	ordinal: string
	index: number
	field: ImportField
}

interface ImportFile {
	name: string
	file: File
	reader: FileReader
	mapping: FieldMap[]
	records: ImportRecord[]
	conflicts: number | null
	preview: [][]
	status: ImportStatus
}

interface ImportRecord {
	barcode: string
	category: string
	model: string
	manufacturer: string
}

let files: ImportFile[] = [];

function parseLine(line: string): string[] {
	let inQuote = false;
	let res = [];
	let final = Array.from(line).reduce((acc, cur) => {
		if (cur === "\"") {
			inQuote = !inQuote;
		} else {
			if (!inQuote && cur === ",") {
				res.push(acc);
				return "";
			} else {
				return acc + cur;
			}
		}
		return acc;
	}, "");
	res.push(final);
	return res;
}

function InventoryImport() {
	const [headersPresent, setHeadersPresent] = useState<boolean>(true);

	const parseFile = (file: File | null, reader: FileReader) => {
		if (file === null || reader.readyState !== 2 || !reader.result || typeof reader.result !== "string") {
			console.error("bad reader/file state");
			return;
		}
		const importFile = files.find((f) => f.file === file);
		if (!importFile) {
			console.log(files);
			console.error("bad react state");
			return;
		}
		let lineEndings = reader.result.match(/\r\n/)
			? "\r\n"
			: reader.result.match(/\r/)
				? "\r" : "\n";
		let lines = reader.result.split(lineEndings);
		// have we already mapped the file? if not, should we auto map it?
		if (importFile.mapping.length === 0) {
			let headers = parseLine((lines.shift() ?? "").toUpperCase());
			console.log(headers);
			if (headersPresent) {
				for (let i = 0; i < headers.length; i++) {
					switch (headers[i]) {
						case "BARCODE":
							if (!importFile.mapping.find((m) => m.field === ImportField.BARCODE)) {
								importFile.mapping.push({
									ordinal: headers[i],
									index: i,
									field: ImportField.BARCODE
								});
							}
							break;
						case "CATEGORY":
							if (!importFile.mapping.find((m) => m.field === ImportField.CATEGORY)) {
								importFile.mapping.push({
									ordinal: headers[i],
									index: i,
									field: ImportField.CATEGORY
								});
							}
							break;
						case "MODEL":
							if (!importFile.mapping.find((m) => m.field === ImportField.MODEL)) {
								importFile.mapping.push({
									ordinal: headers[i],
									index: i,
									field: ImportField.MODEL
								});
							}
							break;
						case "MANUFACTURER":
							if (!importFile.mapping.find((m) => m.field === ImportField.MANUFACTURER)) {
								importFile.mapping.push({
									ordinal: headers[i],
									index: i,
									field: ImportField.MANUFACTURER
								});
							}
							break;
						default:
							importFile.mapping.push({
								ordinal: headers[i],
								index: i,
								field: ImportField.UNUSED
							});
					}
				}
			} else {

			}
		}

		const isMapped = (field: ImportField) => importFile.mapping.find(f => f.field === field);

		if (
			!(isMapped(ImportField.BARCODE) && isMapped(ImportField.CATEGORY)
				&& isMapped(ImportField.MODEL) && isMapped(ImportField.MANUFACTURER))) {
			importFile.status = ImportStatus.UNMAPPED;
			return;
		}
	}

	const handleFile = (evt: ChangeEvent<HTMLInputElement>) => {
		if (evt.target.files === null || evt.target.files.length === 0) {
			return;
		}
		const filesList = evt.target.files;
		const readers: [File, FileReader][] = [];
		for (let i = 0; i < filesList.length; i++) {
			let file = filesList.item(i);
			if (file === null) continue;
			if (!(file.type === "text/csv" || file.name.toLowerCase().endsWith(".csv"))) {
				continue;
			}
			let reader = new FileReader();
			reader.onload = (e) => {
				if (e.target !== null) {
					parseFile(file, e.target);
				}
			};
			reader.readAsText(file);
			readers.push([file, reader]);
		}
		if (readers.length === 0) {
			console.error("no readers created");
		}
		files = readers.map((tuple: [File, FileReader]): ImportFile => ({
			name: tuple[0].name,
			file: tuple[0],
			reader: tuple[1],
			mapping: [],
			records: [],
			preview: [],
			conflicts: null,
			status: ImportStatus.READING
		}));
	}

	let preview = files.map((f) => {

	})

	return (
		<Container>
			<Card>
				<CardContent>
					<FormControlLabel
						control={<Switch inputProps={{'aria-label': 'Headers in first row'}} checked={headersPresent}
						                 onChange={(evt) => setHeadersPresent(evt.target.checked)}/>}
						label={"Headers in first row"}/><br/>
					<input
						accept="text/csv,.csv"
						style={{display: "none"}}
						id="import-file-field"
						type="file" onChange={(evt) => handleFile(evt)}/>
					<label htmlFor="import-file-field">
						<Button variant="contained" component="span">Upload</Button>
					</label>
				</CardContent>
			</Card>

		</Container>
	)
}

export default InventoryImport;
