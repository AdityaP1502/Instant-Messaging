package com.ForwarderServer.Server;

import java.io.File;
import java.io.IOException;
import java.nio.file.*;
import java.util.List;

/* Writer and Reader History
 * All method here used an assumption that the
 * file will be very small
 * By small, it not contains data 
 * more than a mega byte
 * If this assumption is wrong, a more robust method is needed */

// There are two files 
public class HistoryReader {
	private final static String HISTORY_DATA_PATH = "src/server/data/history/";

	private static String readFile(Path path) {
		try {
			List<String> lines = Files.readAllLines(path);
			String read = String.join("|", lines); // need a more robust delimitter
			return read;
		} catch (IOException e) {
			return "";
		}
	}

	private static void writeFile(Path path, String content) throws IOException {
		Files.write(path, (content + System.lineSeparator()).getBytes(), StandardOpenOption.CREATE,
				StandardOpenOption.APPEND);
	}

	public static String readUserHistory(String username) throws IOException {
		String userHistoryPath = HISTORY_DATA_PATH + username + ".txt";
		Path path = Paths.get(userHistoryPath);

		return readFile(path);
	}

	public static void writeUserHistory(String username, String sender, String message, String ISOtimestamp) throws IOException {
		String userHistoryPath = HISTORY_DATA_PATH + username + ".txt";
		Path path = Paths.get(userHistoryPath);

		File f = new File(userHistoryPath);
		f.createNewFile();

		String content = sender + "," + ISOtimestamp + "," + "\"" + message + "\"";

		writeFile(path, content);
	}

	public static void deleteUserHistory(String username) {
		String userHistoryPath = HISTORY_DATA_PATH + username + ".txt";
		Path path = Paths.get(userHistoryPath);
		try {
			Files.deleteIfExists(path);
		} catch (IOException e) {
			System.out.println(e.getMessage());
		}
	}
}
