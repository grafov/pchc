const uint8_t     pinBF = PD2;
const uint8_t     PIN_illumination = A3;                    // номер аналогового вывода к которому подключён освещенности

uint16_t    LUM_result;                               // объявляем переменную для хранения результата опроса датчика освещённости
float       brightness;

void setup() {
    tone(pinBF, 2048, 100);                      // Выводим звуковой сигнал с частотой 2048 Гц и длительностью 0,1 сек
    delay(200);                                  // Не выводим звук в течении 0,1 сек (см. ниже)
    Serial.begin(9600);                                 // инициируем передачу данных в монитор последовательного порта на скорости 9600 бит/сек
    tone(pinBF, 2048, 100);
}

void loop() {
  LUM_result = 1024 - analogRead(PIN_illumination);   // опрашиваем датчик освещённости                    // выводим текст
  brightness = exp((LUM_result+200)/100.)/100.; // min 0 max 280
  if (brightness > 280) {
     brightness = 280;
  }
  Serial.println(brightness);                         // выводим значение освещённости
  delay(3000);
}
