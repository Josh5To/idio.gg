import './App.css';

function test() {
  fetch("http://localhost:3000/.well-known/nostr.json?name=josh")
}

function App() {
  return (
    <div className="App">
      <main>
          <p>nostr pubkey:</p>
          <p>npub1833j3t9d0cc4lu36sj4l8gy5sf7hqkejmpk8tggnpmsqzn67z8nqewh2j8</p>
      </main>
    </div>
  );
}

export default App;
