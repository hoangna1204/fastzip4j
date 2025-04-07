package io.github.hoangna1204.fastzip4j;

import com.sun.jna.Library;

interface FastZip4jLib extends Library {
    void ArchiveFile(String sourceFile, String zipDestination, String temporaryPath, int compressionLevel);

    void ArchiveDir(String sourceDir, String zipFile, String temporaryPath, int compressionLevel);

    void Extract(String zipFile, String destinationDirectory);
}
