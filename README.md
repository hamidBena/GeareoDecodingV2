# Geareo Save Management & Autobuilder Tool

A utility for organizing, modifying, and generating circuits within **Geareo** project files. This tool lets you format your saves for readability, transfer circuits between files, and automatically generate advanced components like custom-sized RGB displays and ROMs.

## Key Features

* **Save File Formatting:** Cleans up and formats raw save data to make it human-readable.
* **Circuit Management:** List all circuits in a save, delete them by name, or copy/move them between different project files.
* **Export & Import:** Share individual circuits with other players via standalone files.
* **Circuit Autobuilder Library:** Automatically generate complex circuits with customizable sizes, including:
  * **RGB Displays** (Custom dimensions)
  * **ROM Circuits** (Custom memory sizes)

## Safety & Backups

Your save data is safe. The tool features an **automatic backup system**:
* Before any save file is overwritten, a backup copy is securely stored.
* If something goes wrong, you can use the built-in **Restore** function to instantly recover your original project file.

## Disclaimer & Liability Warning

**IMPORTANT: USE AT YOUR OWN RISK.**
* **No Liability:** I am not responsible for any data loss, corrupted files, broken circuits, or system errors resulting from the use or misuse of this tool. 
* **Manual Backups:** Even though the tool has an automatic backup and restore engine, it is always a good practice to manually back up your `projects` folder before executing heavy changes.

## Note on Windows Defender / SmartScreen

Because this tool is an independent executable written in Go, Windows may show a "SmartScreen Warning" or flag it as a false positive. 

This happens entirely because the file is new and unsigned. The complete source code is public here for you to audit. To run the tool:
1. Click **More Info** on the Windows pop-up.
2. Click **Run Anyway**.

## Save File Location

Geareo saves your project files locally on your computer here:
`C:\Users\<Your_Username>\AppData\LocalLow\WitWeld\Geareo\projects`

## How to Install and Use

1. Go to the **Releases** tab on the right side of this GitHub page.
2. Download the latest compiled `.exe` application.
3. Launch the app on your Windows PC.
4. Select the function you want to use from the menu.
5. Follow the step-by-step on-screen terminal instructions.

## Bug Reports & Support

If you encounter bugs, broken saves, or want to suggest new layouts for the Autobuilder library, please reach out:
* **Discord DMs:** Message me directly at **@nerdy_human**
* **Discord Official Server:** Send the report directly at the tool's channel: #creations/Geareo save tool

## Credits & License

* Developed by **HamidBn**
* Made for the **Geareo** community.
* Distributed under the MIT License (see the `LICENSE` file for details).
