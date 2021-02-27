# Bienvenido

Bienvenido al proyecto Fan-out / Fan-in.

## Objetivo
El objetivo de este proyecto es mostrarte un ejemplo del patrón Fan-out / Fan-in en Golang para acelerar tus aplicaciones.

## Indice
1. Requisitos.
2. Detalle del patrón.
3. Cómo ejecutar.

### 1. Requisitos.
* Tener instalado [Golang](https://golang.org/dl/) 1.15.* ó superior.
* La velocidad de ejecución dependerá del número de cores presentes en tu CPU, OS y lo que tengas instalado en el. 

### 2. Detalle del patrón.
* Fan-out / Fan-in se refiere al patrón de ejecutar múltiples funciones simultáneamente y luego realizar alguna agregación en los resultados.
![golang workers](./resources/workers.png)
  Atribución de la imagen para Vincent Blanchon

### 3. Cómo ejecutar.
* Desde la línea de comandos ingresa el siguiente comando:
  ````
  /> <directorio>/fan-out-fan-in/go run main.go
 ...