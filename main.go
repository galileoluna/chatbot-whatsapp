package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/mdp/qrterminal"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func GetEventHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			// Obtener el contenido del mensaje
			var messageBody = v.Message.GetConversation()
			// Imprimir el mensaje recibido en consola (opcional)
			fmt.Printf("Mensaje recibido: %s\n", messageBody)

			// Enviar una respuesta (eco del mensaje recibido en este ejemplo)
			client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
				Conversation: proto.String("Recibido: " + messageBody),
			})
		}
	}
}

func main() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	dsn := "user=galibot password=your_secure_password dbname=galibot host=db port=5432 sslmode=disable"
	container, err := sqlstore.New("postgres", dsn, dbLog)
	if err != nil {
		panic(err)
	}

	// Obtener el dispositivo principal
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(GetEventHandler(client))

	if client.Store.ID == nil {
		// Sin ID almacenado, iniciar sesión nuevo
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Renderizar el código QR
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Evento de inicio de sesión:", evt.Event)
			}
		}
	} else {
		// Ya ha iniciado sesión, solo conectar
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Escuchar Ctrl+C para cerrar la aplicación de manera segura
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
