package io.github.hoangna1204.fastzip4j;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;

public class FastZip4jTest {
    
    @Test
    void testCompleteZipWorkflow(@TempDir Path tempDir) throws IOException {
        // Create source directory
        Path sourceDir = tempDir.resolve("source");
        Files.createDirectory(sourceDir);
        
        // Create 10 text files with random content
        for (int i = 0; i < 10; i++) {
            Path file = sourceDir.resolve("file" + i + ".txt");
            String content = "Random content " + UUID.randomUUID();
            Files.writeString(file, content);
        }
        
        // Create target zip file
        Path zipFile = tempDir.resolve("archive.zip");
        
        // Archive the directory with 10 files
        assertDoesNotThrow(() -> {
            FastZip4j.archiveDir(sourceDir, zipFile, 6);
        });
        
        // Create an additional file
        Path additionalFile = tempDir.resolve("additional.txt");
        String additionalContent = "Additional content " + UUID.randomUUID();
        Files.writeString(additionalFile, additionalContent);
        
        // Archive the additional file into the same zip
        assertDoesNotThrow(() -> {
            FastZip4j.archiveFile(additionalFile, zipFile, 6);
        });
        
        // Extract and verify
        Path extractDir = tempDir.resolve("extracted");
        assertDoesNotThrow(() -> {
            FastZip4j.extract(zipFile, extractDir);
        });
        
        // Count the number of files in the extracted directory
        long fileCount = Files.walk(extractDir)
            .filter(Files::isRegularFile)   
            .count();
        
        assertEquals(11, fileCount, "Extracted directory should contain 11 files");
    }

    @Test
    void testArchiveFileWithNonExistentSource(@TempDir Path tempDir) {
        Path nonExistentFile = tempDir.resolve("non_existent_file.txt");
        Path zipFile = tempDir.resolve("archive.zip");
        
        Exception exception = assertThrows(FileNotFoundException.class, () -> {
            FastZip4j.archiveFile(nonExistentFile, zipFile, 6);
        });
        assertTrue(exception.getMessage().toLowerCase().contains("no such file or directory"), 
            "Error message should indicate file not found");
    }

    @Test
    void testExtractToNonExistentTarget(@TempDir Path tempDir) throws IOException, InvalidCompressionLevel {
        // Create a simple zip file first
        Path sourceFile = tempDir.resolve("test.txt");
        String content = "Test content";
        Files.writeString(sourceFile, content);
        
        Path zipFile = tempDir.resolve("archive.zip");
        FastZip4j.archiveFile(sourceFile, zipFile, 6);
        
        // Try to extract to a non-existent directory
        Path nonExistentTarget = tempDir.resolve("extract_target");
        
        // First verify the target doesn't exist
        assertFalse(Files.exists(nonExistentTarget), "Target directory should not exist initially");
        
        // Extract should create the target directory
        FastZip4j.extract(zipFile, nonExistentTarget);
        
        // Verify the extraction
        assertTrue(Files.exists(nonExistentTarget), "Target directory should be created");
        Path extractedFile = nonExistentTarget.resolve("test.txt");
        assertTrue(Files.exists(extractedFile), "Extracted file should exist");
        assertEquals(content, Files.readString(extractedFile), "Extracted content should match original");
    }

    @Test
    void testArchiveDirWithNonExistentSource(@TempDir Path tempDir) {
        Path nonExistentDir = tempDir.resolve("non_existent_dir");
        Path zipFile = tempDir.resolve("archive.zip");
        
        // First create a file to archive so we have a valid zip file
        Path sourceFile = tempDir.resolve("test.txt");
        try {
            Files.writeString(sourceFile, "Test content");
            FastZip4j.archiveFile(sourceFile, zipFile, 6);
        } catch (IOException | InvalidCompressionLevel e) {
            fail("Failed to setup test: " + e.getMessage());
        }
        
        // Now try to archive a non-existent directory into the same zip
        Exception exception = assertThrows(FileNotFoundException.class, () -> {
            FastZip4j.archiveDir(nonExistentDir, zipFile, 6);
        });
        assertTrue(exception.getMessage().toLowerCase().contains("no such file or directory"), 
            "Error message should indicate file not found");
    }
}
