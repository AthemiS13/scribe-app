import React, { useState, useEffect, useRef } from 'react';
import productImage from './assets/images/scribe-display.png';
import textIcon from './assets/icons/text.svg';
import homeIcon from './assets/icons/home.svg';
import settingsIcon from './assets/icons/settings.svg';
import githubIcon from './assets/icons/github.svg';
import uploadIcon from './assets/icons/upload.svg';
import boldIcon from './assets/icons/format-icons/bold.svg';
import italicIcon from './assets/icons/format-icons/italic.svg';
import sizeUpIcon from './assets/icons/format-icons/sizeup.svg';
import sizeDownIcon from './assets/icons/format-icons/sizedown.svg';
import './App.css';              // ← add this
import { Connect, Disconnect, SendData } from "../wailsjs/go/main/App"
import saveIcon from './assets/icons/save.svg';
import importIcon from './assets/icons/import.svg';

// Update the TEXT_CONTAINER to use relative percentages
const TEXT_CONTAINER = {
  widthPercent: 60, // Percent of image width
  heightPercent: 40, // Percent of image height
  lineHeight: 40,
  maxLines: 4,
  padding: 20,
};

const styles = {
  container: {
    width: '100vw',
    height: '100vh',
    backgroundColor: 'black',
    color: 'white',
    display: 'flex',
    flexDirection: 'column',
    fontFamily: 'sans-serif',
    position: 'relative',
    overflow: 'hidden', // Disable scrollingMilitary. Signature, precisely. So you heard how stuff he is. See the bassist sound. Any reason? Yes. Navjot Sidhu wants to. I. So. You do all the time. Something as I'm getting in the country. 
    userSelect: 'text', // Disable text selection/copying
  },
  header: {
    padding: '0px 0 0 0 ', // Increased padding to make header taller
    textAlign: 'center',
    fontSize: '4rem',
    fontFamily: 'Instrument Sans, sans-serif',
    fontWeight: 700,
    width: '100%',
    position: 'relative',
    top: 0,
    height: '120px', // Set fixed height
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: 'white',
  },
  horizontalLine: {
    position: 'absolute',
    left: 0,
    right: 0,
    height: '1px',
    backgroundColor: '#333',
    top: '120px', // Match header height
    margin: 0,
    zIndex: 1,
  },
  mainContent: {
    display: 'flex',
    flex: 1,
  },
  sidebar: {
    width: '120px', // Match vertical line width
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    height: 'calc(100% - 120px)', // Adjust for new header height
    position: 'relative',
    paddingTop: '150px',
  },
  sidebarIconsContainer: {
    display: 'flex',
    flexDirection: 'column',
    gap: '100px', // Increased gap between top icons
    alignItems: 'center',
  },
  sidebarIcon: {
    width: '34px', // Increased by 20%
    height: '34px', // Increased by 20%
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
    transition: 'transform 0.2s ease', // Add transition for hover effect
  },
  githubIcon: {
    position: 'absolute',
    bottom: '70px',
    left: '50%', // Center horizontally
    transform: 'translateX(-50%)', // Center horizontally
    width: '34px', // Increased by 20%
    height: '34px', // Increased by 20%
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
  },
  verticalLine: {
    position: 'absolute',
    top: '0px', // Start from new header height
    right: '-1px',
    width: '1px',
    height: 'calc(100vh - 120px)', // Adjust height based on new header
    backgroundColor: '#333',
    margin: 0,
  },
  textDisplay: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    padding: '0px',
  },
  displayContainer: {
    position: 'relative',
    width: '100%',
    height: 'calc(100vh - 300px)',
    minHeight: '400px',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-end',
    margin: 0,
    paddingRight: '0px', // Add padding for right alignment
    gap: '20px',
    zIndex: 1,
  },
  // Modify styles to use relative positioning
  imageContainer: {
    position: 'relative',
    width: '80%', // Match productImage width
    maxWidth: '2000px', // Match productImage maxWidth
    alignSelf: 'flex-end',
    marginTop: '5%', // Add this line to move the container down
  },
  productImage: {
    width: '100%', // Take full width of container
    height: 'auto',
    display: 'block', // Removes extra space below image
  },
  // Modify styles to use relative positioning
  screenOverlay: {
    position: 'absolute',
    top: '47.8%', // Adjust these percentages to position correctly
    left: '46.3%',
    transform: 'translate(-50%, -50%)',
    width: '53%', // Relative to image width
    height: '72%', // Relative to image height
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center', // Center content vertically
    color: 'white',
    //padding: `${TEXT_CONTAINER.padding}px`,
    zIndex: 2,
    pointerEvents: 'none',
    backgroundColor: 'transparent',
  },
  pageText: {
    width: '100%',
    height: '100%',
    whiteSpace: 'pre-wrap',
    wordWrap: 'break-word',
    overflow: 'hidden',
    textAlign: 'left',
    // Font size will be set dynamically
  },
  // Modify styles to use relative positioning
  navigationOverlay: {
    position: 'absolute',
    top: '50%', // Changed for better vertical centering
    transform: 'translateY(-50%)', // Added for better vertical centering
    left: '0',
    right: '0',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '0 18% 0 10%', // Adjusted right padding to bring arrows closer
    zIndex: 3,
    pointerEvents: 'none',
  },
  arrow: {
    fontSize: '4vw', // Use vw for responsive scaling with viewport width
    color: 'black',
    cursor: 'pointer',
    userSelect: 'none',
    pointerEvents: 'auto', // Re-enable clicks for arrows
  },
  inputSection: {
    width: '80%',
    display: 'flex',
    flexDirection: 'column', // Stack elements vertically
    gap: '15px',
    alignItems: 'center',
    marginTop: '100px', // Move up
    position: 'relative',
    zIndex: 4, // Ensure input section is above overlays
  },
  inputRow: {
    display: 'flex',
    gap: '20px',
    alignItems: 'center',
    width: '100%',
  },
  input: {
    flex: 1,
    padding: '15px',
    borderRadius: '25px',
    backgroundColor: '#111',
    border: '1px solid #333',
    color: 'white',
    fontSize: '16px',
    userSelect: 'text', // Enable text selection
    cursor: 'text', // Show text cursor
  },
  formatButtons: {
    display: 'flex',
    gap: '20px',
    alignSelf: 'flex-start', // Align with input field
    marginLeft: '0', // Align with input field
    marginTop: '1%', // Add some space above the buttons
  },
  formatButton: {
    width: '64px', // Increased size
    height: '64px', // Increased size
    borderRadius: '15px',
    backgroundColor: '#111',
    border: '1px solid #333',
    padding: '20px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
    transition: 'transform 0.2s ease', // Add transition for hover effect
    '&:hover': {
      transform: 'scale(1.1)', // 10% size increase on hover
    }
  },
  sendButton: {
    padding: '15px 30px',
    borderRadius: '25px',
    backgroundColor: '#111',
    border: '1px solid #333',
    color: 'white',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    gap: '10px',
  },
  uploadIconStyle: {
    width: '16px', // Smaller upload icon
    height: '16px',
    padding: 0,
  },
};

const iconStyle = {
  width: '100%',
  height: '100%',
  padding: '2px',
};

// Add this constant for the max line length
const MAX_LINE_LENGTH = 21;
const MANUAL_BREAK_MARKER = '\u200B'; // Zero-width space (invisible, safe for input)

function App() {
  const [inputText, setInputText] = useState('');
  const [fontSize, setFontSize] = useState(24); 
  const [fontSizeMultiplier, setFSMultiplier] = useState(1);
  const [pages, setPages] = useState([[]]);
  const [currentPage, setCurrentPage] = useState(0);
  // Add new state for the progress popup
  const [showProgress, setShowProgress] = useState(false);
  const [progress, setProgress] = useState(0);
  const [uploadComplete, setUploadComplete] = useState(false);
  const [pageInputs, setPageInputs] = useState(['']); // Store raw input text for each page
  const [manualBreaks, setManualBreaks] = useState([[]]); // Array of arrays of break indices
  const displayAreaRef = useRef(null)

  useEffect(() => {
    if (!displayAreaRef) return;

    const resizeObserver = new ResizeObserver((entries) => {
      for (let entry of entries) {
        if (entry.contentBoxSize) {
          setFontSize(entry.contentRect.width / 13)
        }
      }
    });

    resizeObserver.observe(displayAreaRef.current);

    return () => {
      resizeObserver.disconnect();
    };
  }, [])

  // Break text at container boundaries, grouping into pages of max 4 lines
  const formatText = (text) => {
    if (!text) return [['Your Text Will Appear Here!']];
    
    const paragraphs = text.split('\n');
    const lines = [];
    
    const measureEl = document.createElement('div');
    measureEl.style.position = 'absolute';
    measureEl.style.visibility = 'hidden';
    // Ensure fontSize is applied correctly, potentially using a state variable if it changes
    // For example, if you have a state `currentFontSize`
    measureEl.style.fontSize = `${fontSize * fontSizeMultiplier}px`; 
    measureEl.style.fontFamily = 'sans-serif'; // Or your specific preview font
    measureEl.style.whiteSpace = 'nowrap';
    document.body.appendChild(measureEl);
    
    // Calculate maxWidth based on the displayAreaRef (preview area)
    // and TEXT_CONTAINER.widthPercent
    let maxWidth = 300; // Default or fallback width
    if (displayAreaRef.current) {
      const previewContainerWidth = displayAreaRef.current.offsetWidth;
      maxWidth = (previewContainerWidth * (TEXT_CONTAINER.widthPercent / 100)) - (TEXT_CONTAINER.padding * 2);
    }
    
    paragraphs.forEach(paragraph => {
      if (paragraph === '') {
        lines.push(''); // Add an empty line for explicit newlines
        return;
      }
      
      const words = paragraph.split(' ');
      let currentLine = '';
      
      for (const word of words) {
        const testLine = currentLine ? `${currentLine} ${word}` : word;
        measureEl.textContent = testLine;
        if (measureEl.offsetWidth > maxWidth && currentLine) {
          lines.push(currentLine);
          currentLine = word;
        } else {
          currentLine = testLine;
        }
      }
      if (currentLine) {
        lines.push(currentLine);
      }
    });
    
    document.body.removeChild(measureEl);
    
    const result = [];
    for (let i = 0; i < lines.length; i += TEXT_CONTAINER.maxLines) {
      result.push(lines.slice(i, i + TEXT_CONTAINER.maxLines));
    }
    return result.length ? result : [['Your Text Will Appear Here!']];
  };

  // Recompute pages when inputText or fontSize changes
  useEffect(() => {
    // Save the current text to pageInputs first
    const newPageInputs = [...pageInputs];
    newPageInputs[currentPage] = inputText;
    setPageInputs(newPageInputs);
    
    // Update only the current page in the pages array
    const newPages = [...pages];
    newPages[currentPage] = formatText(inputText)[0]; // Format just the current page text
    setPages(newPages);
  }, [inputText, fontSize]);

  const increaseFontSize = () => {
    setFSMultiplier((prev) => prev === 3 ? prev : prev + 1)
  };

  const decreaseFontSize = () => {
    setFSMultiplier((prev) => prev === 1 ? 1 : prev - 1)
  };

  // Update the navigation handlers
  const handlePrevPage = () => {
    // Save current page text first
    const newPageInputs = [...pageInputs];
    newPageInputs[currentPage] = inputText;
    setPageInputs(newPageInputs);
    
    // Navigate to previous page with cycling
    const prevPage = currentPage > 0 ? currentPage - 1 : pages.length - 1;
    setCurrentPage(prevPage);
    
    // Update input field with previous page's content
    setInputText(pageInputs[prevPage] || '');
  };

  const handleNextPage = () => {
    // Save current page content first
    const newPageInputs = [...pageInputs];
    newPageInputs[currentPage] = inputText;
    setPageInputs(newPageInputs);
    
    // Check if we're on the last page
    if (currentPage === pages.length - 1) {
      // Check if current page has actual text content
      const hasContent = pages[currentPage][0] !== 'Your Text Will Appear Here!' && 
                         pages[currentPage].length > 0 && 
                         pages[currentPage][0] !== '';
      
      if (hasContent) {
        // Clone the pages array and add a new blank page
        const newPages = [...pages];
        newPages.push(['']);
        setPages(newPages);
        
        // Add new empty entry to pageInputs
        const updatedPageInputs = [...newPageInputs];
        updatedPageInputs.push('');
        setPageInputs(updatedPageInputs);
        
        // Move to the new page
        const nextPage = currentPage + 1;
        setCurrentPage(nextPage);
        
        // Clear input field for the new blank page
        setInputText('');
      } else {
        // If on last page with no content, cycle to first page
        setCurrentPage(0);
        setInputText(pageInputs[0] || '');
      }
    } else {
      // Just navigate to the next existing page
      const nextPage = currentPage + 1;
      setCurrentPage(nextPage);
      setInputText(pageInputs[nextPage] || '');
    }
  };

  // Helper to insert a manual break at the current cursor position
  const insertManualBreak = () => {
    const input = document.activeElement;
    if (!input || input.tagName !== "INPUT") return;
    const cursorPos = input.selectionStart;
    const before = inputText.slice(0, cursorPos);
    const after = inputText.slice(cursorPos);
    setInputText(before + MANUAL_BREAK_MARKER + after);
    // Move cursor after the marker (optional, but recommended)
    setTimeout(() => {
      input.setSelectionRange(cursorPos + 1, cursorPos + 1);
    }, 0);
  };

  // Handle Ctrl+Enter for manual break
  const handleInputKeyDown = (e) => {
    if (e.ctrlKey && e.key === "Enter") {
      e.preventDefault();
      insertManualBreak();
    }
  };

  // Format text for preview and sending (applies breaks)
  const getFormattedLines = (text) => {
    if (!text) return ['Your Text Will Appear Here!'];
    let lines = [];
    let segments = text.split(MANUAL_BREAK_MARKER);

    for (let segment of segments) {
      while (segment.length > MAX_LINE_LENGTH) {
        lines.push(segment.slice(0, MAX_LINE_LENGTH));
        segment = segment.slice(MAX_LINE_LENGTH);
      }
      lines.push(segment);
    }
    // Ensure 4 lines per page
    while (lines.length < 4) lines.push('');
    return lines.slice(0, 4);
  };

  // New helper to get all formatted lines from all pages
  const getAllFormattedLines = () => {
    const allLines = [];
    for (let i = 0; i < pages.length; i++) {
      const pageLines = getFormattedLines(pageInputs[i]);
      allLines.push(...pageLines);
    }
    return allLines;
  };

  async function send() {
    const connected = await Connect();
    if (connected) {
      // Get all lines from all pages, each page is 4 lines
      const allLines = getAllFormattedLines();
      // Join with '\n' to ensure a break after every line (including after every 4th line)
      const sendString = allLines.join('\n');
      const sent = await SendData(sendString);
      if (!sent) {
        console.log("pruser");
      }
      await Disconnect();
    }
  }

  // Handle the send button click
  const handleSend = () => {
    send()
    setProgress(0);
    setUploadComplete(false);
    setShowProgress(true);

    // Simulate progress
    let currentProgress = 0;
    const interval = setInterval(() => {
      currentProgress += 5;
      setProgress(currentProgress);

      if (currentProgress >= 100) {
        clearInterval(interval);
        setUploadComplete(true);

        // Auto-close after 1 second
        setTimeout(() => {
          setShowProgress(false);
        }, 1000);
      }
    }, 100); // Update every 100ms
  };

  // Handle input changes
  const handleInputChange = (e) => {
    // Replace all '-' with the marker for internal storage
    let newText = e.target.value.replace(/-/g, MANUAL_BREAK_MARKER);

    // Split into segments by manual break marker
    let segments = newText.split(MANUAL_BREAK_MARKER);
    let lines = [];
    for (let segment of segments) {
      while (segment.length > MAX_LINE_LENGTH) {
        lines.push(segment.slice(0, MAX_LINE_LENGTH));
        segment = segment.slice(MAX_LINE_LENGTH);
      }
      lines.push(segment);
    }

    // If more than 4 lines, trim extra
    if (lines.length > 4) {
      lines = lines.slice(0, 4);
      // Rebuild text with markers
      let rebuilt = '';
      let charCount = 0;
      for (let i = 0; i < segments.length && charCount < 4; i++) {
        let seg = segments[i];
        while (seg.length > MAX_LINE_LENGTH && charCount < 4) {
          rebuilt += seg.slice(0, MAX_LINE_LENGTH) + MANUAL_BREAK_MARKER;
          seg = seg.slice(MAX_LINE_LENGTH);
          charCount++;
        }
        if (charCount < 4) {
          rebuilt += seg;
          charCount++;
          if (i < segments.length - 1 && charCount < 4) rebuilt += MANUAL_BREAK_MARKER;
        }
      }
      newText = rebuilt;
    } else if (lines.length === 4 && lines[3].length > MAX_LINE_LENGTH) {
      // Prevent typing more than 21 chars on the 4th line
      lines[3] = lines[3].slice(0, MAX_LINE_LENGTH);
      newText = lines.slice(0, 4).join('');
      // Add back manual breaks
      let rebuilt = '';
      let idx = 0;
      for (let s of segments) {
        if (idx >= 4) break;
        let chunk = s.slice(0, MAX_LINE_LENGTH * (4 - idx));
        rebuilt += chunk;
        idx += Math.ceil(chunk.length / MAX_LINE_LENGTH) || 1;
        if (idx < 4) rebuilt += MANUAL_BREAK_MARKER;
      }
      newText = rebuilt;
    }

    setInputText(newText);
  };

  // Add file handling functions
  const handleSaveFile = () => {
    try {
      // Create a JSON representation of all pages with their breaks
      const dataToSave = {
        pages: pageInputs,
        version: "1.0"
      };
      
      const jsonString = JSON.stringify(dataToSave, null, 2);
      const blob = new Blob([jsonString], { type: 'application/json' });
      
      // Create a download link and trigger it
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'scribe-content.json';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error("Failed to save file:", error);
      // You could add a user-friendly error message here
    }
  };
  
  const handleImportFile = () => {
    try {
      // Create a file input element and trigger it
      const fileInput = document.createElement('input');
      fileInput.type = 'file';
      fileInput.accept = '.json';
      
      fileInput.onchange = (e) => {
        const file = e.target.files[0];
        if (!file) return;
        
        const reader = new FileReader();
        reader.onload = (event) => {
          try {
            const importedData = JSON.parse(event.target.result);
            
            // Validate the imported data has the expected structure
            if (importedData && Array.isArray(importedData.pages)) {
              // Update the state with imported pages
              setPageInputs(importedData.pages);
              
              // Generate pages data from imported page inputs
              const newPages = importedData.pages.map(pageText => 
                formatText(pageText)[0]
              );
              setPages(newPages);
              
              // Reset to first page and update input text
              setCurrentPage(0);
              setInputText(importedData.pages[0] || '');
            } else {
              console.error("Invalid file format");
              // You could add a user-friendly error message here
            }
          } catch (parseError) {
            console.error("Failed to parse file:", parseError);
            // You could add a user-friendly error message here
          }
        };
        reader.readAsText(file);
      };
      
      fileInput.click();
    } catch (error) {
      console.error("Failed to import file:", error);
      // You could add a user-friendly error message here
    }
  };

  return (
    <div style={styles.container}>
      <header style={styles.header} className="header-text">scribe</header>
      <div style={styles.horizontalLine} />
      <div style={styles.mainContent}>
        <div style={styles.sidebar}>
          <div style={styles.sidebarIconsContainer}>
            <div style={styles.sidebarIcon} className="sidebar-icon">
              <img src={textIcon} alt="Text" style={iconStyle} />
            </div>
            <div style={styles.sidebarIcon} className="sidebar-icon">  {/* add className */}
              <img src={homeIcon} alt="Home" style={iconStyle} />
            </div>
            <div style={styles.sidebarIcon} className="sidebar-icon">  {/* add className */}
              <img src={settingsIcon} alt="Settings" style={iconStyle} />
            </div>
          </div>
          <a href="https://github.com/AthemiS13/scribe-app" style={styles.githubIcon}>
            <img src={githubIcon} alt="GitHub" style={iconStyle} />
          </a>
          <div style={styles.verticalLine} />
        </div>
        <div style={styles.textDisplay}>
          <div style={styles.imageContainer}>
            <img
              src={productImage}
              alt="Product Display"
              style={styles.productImage}
            />
            <div ref={displayAreaRef} style={styles.screenOverlay}>
              <div
                style={{
                  ...styles.pageText,
                  fontSize: `${fontSize * fontSizeMultiplier}px`,
                }}
              >
                {getFormattedLines(pageInputs[currentPage]).map((line, idx) => (
                  <div key={idx}>{line}</div>
                ))}
              </div>
            </div>
            <div style={styles.navigationOverlay}>
              <span style={styles.arrow} onClick={handlePrevPage}>&lt;</span>
              <span style={styles.arrow} onClick={handleNextPage}>&gt;</span>
            </div>
          </div>

          {/* Input section remains the same */}
          <div style={styles.inputSection}>
            <div style={styles.inputRow}>
              <input
                type="text"
                placeholder="Type Here!"
                style={styles.input}
                value={inputText.replace(/\u200B/g, '-')}
                onChange={handleInputChange}
                onKeyDown={handleInputKeyDown}
              />
              <button
                className="send-button"
                onClick={handleSend}
              >
                <span>Send</span>
                <img src={uploadIcon} alt="Upload" style={styles.uploadIconStyle} />
              </button>
            </div>
            <div style={styles.formatButtons}>
              <button
                style={styles.formatButton}
                onClick={increaseFontSize}
                className="format-button"
              >
                <img src={sizeUpIcon} alt="Size Up" style={{ width: '100%', height: '100%', padding: '2px' }} />
              </button>
              <button
                style={styles.formatButton}
                onClick={decreaseFontSize}
                className="format-button"
              >
                <img src={sizeDownIcon} alt="Size Down" style={{ width: '100%', height: '100%', padding: '2px' }} />
              </button>
              <button style={styles.formatButton} className="format-button">
                <img src={boldIcon} alt="Bold" style={{ width: '100%', height: '100%', padding: '2px' }} />
              </button>
              <button
                style={styles.formatButton}
                className="format-button"
                onClick={insertManualBreak}
                type="button"
                title="Insert Manual Break"
              >
                <img src={italicIcon} alt="Manual Break" style={{ width: '100%', height: '100%', padding: '2px' }} />
              </button>
              
              {/* Gap between format buttons and file operation buttons */}
              <div style={{ width: '60px' }}></div>
              
              {/* New Save button */}
              <button
                style={styles.formatButton}
                className="format-button"
                onClick={handleSaveFile}
                type="button"
                title="Save to File"
              >
                <img src={saveIcon} alt="Save" style={{ width: '100%', height: '100%', padding: '2px' }} />
              </button>
              
              {/* New Import button */}
              <button
                style={styles.formatButton}
                className="format-button"
                onClick={handleImportFile}
                type="button"
                title="Import from File"
              >
                <img src={importIcon} alt="Import" style={{ width: '100%', height: '100%', padding: '2px' }} />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Progress Popup */}
      {showProgress && (
        <div className="progress-overlay">
          <div className="progress-container">
            <div
              className="progress-circle"
              style={{ "--progress": `${progress}%` }}
            >
              <div className="progress-text">
                {uploadComplete ? "Done!" : `${progress}%`}
              </div>
            </div>
            <div className="progress-label">
              {uploadComplete ? "Complete" : "Uploading"}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
