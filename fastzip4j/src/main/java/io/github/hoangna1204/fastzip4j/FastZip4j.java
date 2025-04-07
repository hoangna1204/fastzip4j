package io.github.hoangna1204.fastzip4j;

import com.sun.jna.Native;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Comparator;

/**
 * The {@code FastZip4j} class provides utility methods for archiving files and
 * directories
 * into ZIP files and extracting ZIP files to a specified directory. It uses
 * native library
 * bindings for fast and efficient compression and decompression.
 */
public class FastZip4j {
    private static final FastZip4jLib fastzip4jLib;

    static {
        var osName = System.getProperty("os.name").toLowerCase();
        var arch = System.getProperty("os.arch").toLowerCase();
        var libraryName = "fastzip4j";
        if (osName.contains("mac")) {
            if (arch.contains("x86_64") || arch.contains("amd64")) {
                libraryName = "fastzip4j_amd64";
            } else if (arch.contains("arm") || arch.contains("aarch64")) {
                libraryName = "fastzip4j_arm64";
            } else {
                throw new UnsupportedOperationException("Unsupported architecture: " + arch);
            }
        }
        try {
            fastzip4jLib = Native.load(libraryName, FastZip4jLib.class);
        } catch (UnsatisfiedLinkError e) {
            throw new LibraryLoadFailedException("Failed to load native library: " + e.getMessage());
        }

    }

    private FastZip4j() {
    }

    /**
     * Archives a single file into a ZIP file with the specified compression level.
     *
     * @param sourceFile       the file to be archived
     * @param zipFile          the destination ZIP file
     * @param compressionLevel the compression level. CompressionLevel must be
     *                         between 1 (BestSpeed) and 9 (BestCompression). Higher
     *                         levels typically run slower but compress more.
     * @throws InvalidCompressionLevel if the compression level is not between 1 and
     *                                 9
     * @throws IOException             if the source file or directory does not
     *                                 exist
     */
    public static void archiveFile(Path sourceFile, Path zipFile, int compressionLevel)
            throws InvalidCompressionLevel, IOException {
        if (!Files.exists(sourceFile)) {
            throw new FileNotFoundException("No such file or directory: " + sourceFile);
        }
        if (!Files.isRegularFile(sourceFile)) {
            throw new FileNotFoundException("No such file or directory: " + sourceFile);
        }
        if (compressionLevel < 1 || compressionLevel > 9) {
            throw new InvalidCompressionLevel("Compression level must be between 1 and 9");
        }
        var tempDir = Files.createTempDirectory("fastzip4j");
        archiveFile(sourceFile.toAbsolutePath().normalize().toString(), zipFile.toAbsolutePath().normalize().toString(),
                tempDir.toAbsolutePath().normalize().toString(), compressionLevel);
        deleteDir(tempDir);
    }

    /**
     * Archives a directory and its contents into a ZIP file with the specified
     * compression level.
     *
     * @param sourceDir        the directory to be archived
     * @param zipFile          the destination ZIP file
     * @param compressionLevel the compression level. CompressionLevel must be
     *                         between 1 (BestSpeed) and 9 (BestCompression). Higher
     *                         levels typically run slower but compress more.
     * @throws InvalidCompressionLevel
     * @throws IOException             if the source directory or file does not
     *                                 exist
     */
    public static void archiveDir(Path sourceDir, Path zipFile, int compressionLevel)
            throws InvalidCompressionLevel, IOException {
        if (!Files.exists(sourceDir)) {
            throw new FileNotFoundException("No such file or directory: " + sourceDir);
        }
        if (!Files.isDirectory(sourceDir)) {
            throw new FileNotFoundException("No such file or directory: " + sourceDir);
        }
        if (compressionLevel < 1 || compressionLevel > 9) {
            throw new InvalidCompressionLevel("Compression level must be between 1 and 9");
        }
        var tempDir = Files.createTempDirectory("fastzip4j");
        archiveDir(sourceDir.toAbsolutePath().normalize().toString(), zipFile.toAbsolutePath().normalize().toString(),
                tempDir.toAbsolutePath().normalize().toString(), compressionLevel);
        deleteDir(tempDir);
    }

    /**
     * Extracts the contents of a ZIP file into the specified destination directory.
     *
     * @param zipFile              the ZIP file to be extracted
     * @param destinationDirectory the destination directory where the contents will
     *                             be extracted
     */
    public static void extract(Path zipFile, Path destinationDirectory) {
        extract(zipFile.toAbsolutePath().normalize().toString(),
                destinationDirectory.toAbsolutePath().normalize().toString());
    }

    private static void archiveFile(String sourceFile, String zipFile, String temporaryPath, int compressionLevel) {
        fastzip4jLib.ArchiveFile(sourceFile, zipFile, temporaryPath, compressionLevel);
    }

    private static void archiveDir(String sourceDir, String zipFile, String temporaryPath, int compressionLevel) {
        fastzip4jLib.ArchiveDir(sourceDir, zipFile, temporaryPath, compressionLevel);
    }

    private static void extract(String zipFile, String destinationDirectory) {
        fastzip4jLib.Extract(zipFile, destinationDirectory);
    }

    private static void deleteDir(Path dir) throws IOException {
        if (!Files.exists(dir)) {
            return;
        }
        try (var paths = Files.walk(dir)) {
            paths.sorted(Comparator.reverseOrder()).map(Path::toFile).forEach(File::delete);
        }
    }
}
