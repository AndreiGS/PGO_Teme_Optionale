import { useState } from "react";
import { Button } from "./button";
import { Progress } from "./progress";
import { EventsOff, EventsOn } from '../../../wailsjs/runtime/runtime'
import { GenerateFromFile } from "../../../wailsjs/go/main/App";

const FileUpload = () => {
    const [selectedFiles, setSelectedFiles] = useState<FileList | null>(null);
    const [generalProgress, setGeneralProgress] = useState<number>(0);
    const [filesProgress, setFilesProgress] = useState<number[]>([]);

    const handleFileChange = (event: any) => {
        setSelectedFiles(event.target.files);
        setGeneralProgress(0);
        setFilesProgress([]);
    };

    const handleFile = async (file: File): Promise<void> => {
        await GenerateFromFile(file.name);
    };

    const handleFilesUpload = async () => {
        if (selectedFiles == null) {
            alert("Please select a file to upload.");
            return;
        }

        setFilesProgress(Array.from({ length: selectedFiles.length }, () => 0));
        for (let index = 0; index < selectedFiles.length; index++) {
            const file = selectedFiles[index];
            EventsOn("fileupload", (result: any) => {
                console.log(result);
                const { noLines, maxLines } = result;
                const newFilesProgress = [...filesProgress];
                newFilesProgress[index] = noLines * 100 / maxLines;
                setFilesProgress(newFilesProgress);
            });
            await handleFile(file);
            EventsOff("fileupload");
            setGeneralProgress((index + 1) * 100 / selectedFiles.length);
        }
    };

    return (
        <div className="flex flex-col gap-4">
            <input type="file" accept=".txt" onChange={handleFileChange} multiple />
            <Button onClick={async () => {
                await handleFilesUpload();
            }}>Generate</Button>
            {
                selectedFiles && (
                    <div>
                        <div>
                            {filesProgress.map((progress, index) => (
                                <div key={index}>
                                    File {index + 1} progress:
                                    <Progress value={progress} />
                                </div>
                            ))}
                        </div>
                        <div>
                            General progress:
                            <Progress value={generalProgress} />
                        </div>
                    </div>
                )
            }
        </div>
    );
};

export default FileUpload;