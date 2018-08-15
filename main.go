package main

import (
  "os"
  "log"
  "fmt"
  "bufio"
  "strings"
  "regexp"
  "syscall"
  "net/url"
  "net/http"
  "io/ioutil"
  "crypto/tls"
  "golang.org/x/crypto/ssh/terminal"
)

func main(){
  // Skip Certificate verification
  http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

  // Regexp Setup
  ReLogin := regexp.MustCompile("You have successfully logged in");
  ReMaxLimit := regexp.MustCompile("You have reached Maximum Login Limit");
  ReInvalidCred := regexp.MustCompile("Make sure your password is correct");

  // Reader
  reader := bufio.NewReader(os.Stdin);

  // Getting Creadentials
  fmt.Print("Username: ");
  inputUsername, _ := reader.ReadString('\n');

  fmt.Print("Password: ");
  bytePassword, err := terminal.ReadPassword(int(syscall.Stdin));
  fmt.Println("\nConnecting..");
  if err != nil {
    log.Fatal(err);
  }

  // Extracted Credentials
  password := string(bytePassword);
  username := strings.Split(inputUsername, "\n")[0];

  resp, err := http.PostForm("https://10.100.56.55:8090/login.xml",
    url.Values{
      "a": {"1524343263066"},
      "mode": {"191"},
      "username": {username},
      "password": {password},
      "producttype": {"0"}});
  if err != nil {
    log.Fatal(err)
  }

  defer resp.Body.Close();

  XMLBytes, err := ioutil.ReadAll(resp.Body);

  if err != nil {
    log.Fatal(err)
  }

  // Checking Status
  if ReLogin.FindSubmatch(XMLBytes) != nil {
    fmt.Println("Succeefully Logged In.");
  } else if ReMaxLimit.FindSubmatch(XMLBytes) != nil {
    fmt.Println("Maximum Login Limit Reach.");
  } else if ReInvalidCred.FindSubmatch(XMLBytes) != nil {
    fmt.Println("Invalid Credentials.");
  } else {
    fmt.Println("Data Exceed.");
  }
}