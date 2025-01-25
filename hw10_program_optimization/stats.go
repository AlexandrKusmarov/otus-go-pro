package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat, err := getDomainCount(r, domain)
	if err != nil {
		return nil, fmt.Errorf("getDomainCount error: %w", err)
	}
	return domainStat, err
}

func getDomainCount(r io.Reader, domain string) (result DomainStat, err error) {
	var parser fastjson.Parser
	finalResult := make(DomainStat)
	if domain == "" {
		return nil, nil
	}
	domain = strings.ToLower(domain)

	// Используем bufio для построчного чтения
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// Парсим JSON
		v, err := parser.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}

		// Получаем значение email
		email := string(v.GetStringBytes("Email"))
		atIndex := strings.LastIndex(email, "@")
		if atIndex == -1 || atIndex == len(email)-1 {
			continue // Пропускаем некорректные email
		}

		emailDomain := strings.ToLower(email[atIndex+1:]) // Получаем домен
		if strings.HasSuffix(emailDomain, "."+domain) {
			finalResult[emailDomain]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return finalResult, nil
}
