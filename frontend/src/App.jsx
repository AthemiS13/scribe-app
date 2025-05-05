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
    overflow: 'hidden', // Disable scrolling
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
    top: '38%', // Match the text overlay
    left: '0',
    right: '0',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '0 20%', // Adjust arrow positioning
    zIndex: 3,
    pointerEvents: 'none',
  },
  arrow: {
    fontSize: '24px',
    color: 'white',
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
    if (!text) return [['Your Text Will Appear Here']];
    const measureEl = document.createElement('div');
    measureEl.style.position = 'absolute';
    measureEl.style.visibility = 'hidden';
    measureEl.style.fontSize = `${fontSize}px`;
    measureEl.style.fontFamily = 'sans-serif';
    measureEl.style.whiteSpace = 'nowrap';
    document.body.appendChild(measureEl);

    const words = text.split(' ');
    const lines = [];
    let currentLine = '';
    const maxWidth = TEXT_CONTAINER.width - TEXT_CONTAINER.padding * 2;

    for (const word of words) {
      // Test if adding this word fits in current line
      measureEl.textContent = currentLine ? `${currentLine} ${word}` : word;
      if (measureEl.offsetWidth > maxWidth) {
        if (currentLine) lines.push(currentLine);
        currentLine = word;
      } else {
        currentLine = currentLine ? `${currentLine} ${word}` : word;
      }
    }
    if (currentLine) lines.push(currentLine);

    document.body.removeChild(measureEl);

    // Group lines into pages (4 lines per page)
    const result = [];
    for (let i = 0; i < lines.length; i += TEXT_CONTAINER.maxLines) {
      result.push(lines.slice(i, i + TEXT_CONTAINER.maxLines));
    }
    return result.length ? result : [['Your Text Will Appear Here']];
  };

  // Recompute pages when inputText or fontSize changes
  useEffect(() => {
    const newPages = formatText(inputText);
    setPages(newPages);
    setCurrentPage(0); // Reset to first page
  }, [inputText, fontSize]);

  const increaseFontSize = () => {
    setFSMultiplier((prev) => prev === 3 ? prev : prev + 1)
  };

  const decreaseFontSize = () => {
    setFSMultiplier((prev) => prev === 1 ? 1 : prev - 1)
  };

  const handlePrevPage = () => {
    setCurrentPage((p) => (p > 0 ? p - 1 : p));
  };

  const handleNextPage = () => {
    setCurrentPage((p) => (p < pages.length - 1 ? p + 1 : p));
  };

  // Handle the send button click
  const handleSend = () => {
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
                {pages[currentPage]?.map((line, idx) => (
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
                value={inputText}
                onChange={(e) => setInputText(e.target.value)}
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
              <button style={styles.formatButton} className="format-button"> {/* add className */}
                <img src={boldIcon} alt="Bold" style={{ width: '100%', height: '100%', padding: '2px' }} />
              </button>
              <button style={styles.formatButton} className="format-button"> {/* add className */}
                <img src={italicIcon} alt="Italic" style={{ width: '100%', height: '100%', padding: '2px' }} />
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
